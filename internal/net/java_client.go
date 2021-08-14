package net

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net/auth"
	. "github.com/JungleMC/protocol"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"reflect"
)

type JavaClient struct {
	networkId           uuid.UUID
	server              *JavaServer
	conn                net.Conn
	reader              *bufio.Reader
	writer              *bufio.Writer
	state               ConnectionState
	gameProtocolVersion int32
	compressionEnabled  bool
	encryptionEnabled   bool
	verifyToken         []byte
	sharedSecret        []byte
	encryptStream       cipher.Stream
	decryptStream       cipher.Stream
	authProfile         *auth.Profile
}

func readAll(r io.Reader) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			return b, err
		}
	}
}

func (c *JavaClient) listen() {
	for {
		pkt, err := c.decodePacket(c.reader)
		if err != nil {
			c.conn.Close()
			return
		}

		if pkt == nil {
			panic("pkt was nil")
		}

		err = c.handle(pkt)
		if err != nil {
			c.disconnectError(err)
		}
	}
}

func (c *JavaClient) Send(pkt Packet) error {
	log.Printf("tx -> %v\n", reflect.ValueOf(pkt).Elem().Type().Name())

	// TODO: Submit packets to a FIFO queue before sending directly, maintaining packet order
	defer c.writer.Flush()
	_, err := c.writer.Write(c.encodePacket(pkt))
	if err != nil {
		return err
	}
	return nil
}

func (c *JavaClient) encodePacket(pkt Packet) []byte {
	buf := WriteVarInt32(pkt.Id())
	buf = append(buf, pkt.Encode()...)

	if c.compressionEnabled {
		if len(buf) >= config.Get.CompressionThreshold {
			buf = compress(buf)
			buf = append(WriteVarInt32(int32(len(buf))), buf...)
		} else {
			buf = append(WriteVarInt32(0), buf...)
		}
	}
	buf = append(WriteVarInt32(int32(len(buf))), buf...)

	if config.Get.OnlineMode && c.encryptionEnabled {
		c.encryptStream.XORKeyStream(buf, buf)
	}
	return buf
}

func readVarInt32(buf *bufio.Reader) (int32, error) {
	numRead := int32(0)
	result := int32(0)

	for {
		read, err := buf.ReadByte()
		if err != nil {
			return 0, err
		}
		value := read & 0b01111111
		result |= int32(value) << (7 * numRead)
		numRead++
		if numRead > 5 {
			panic("varint32 is too big")
		}

		if (read & 0b10000000) == 0 {
			break
		}
	}
	return result, nil
}

func (c *JavaClient) decodePacket(buf *bufio.Reader) (Packet, error) {
	if c.state == Handshake {
		peekByte, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		if peekByte == LegacyPingPayload {
			pkt := NewPacket(c.state, ClientToServer, int32(peekByte))
			pkt.Decode(buf)
			return pkt, nil
		} else {
			buf.UnreadByte()
		}
	}

	dataLen, err := readVarInt32(buf)
	if err != nil {
		return nil, err
	}

	if c.compressionEnabled {
		uncompressedLen := ReadVarInt32(buf)

		if uncompressedLen > 0 {
			compressed := make([]byte, dataLen)
			buf.Read(compressed)

			uncompressed := make([]byte, uncompressedLen)
			zl, _ := zlib.NewReader(bytes.NewBuffer(compressed))
			zl.Read(uncompressed)
			zl.Close()
			buf = bufio.NewReader(bytes.NewReader(uncompressed))
		}
	}

	id := ReadVarInt32(buf)
	pkt := NewPacket(c.state, ClientToServer, id)
	if pkt == nil {
		panic("pkt is nil")
	}
	pkt.Decode(buf)

	log.Printf("rx <- %v\n", reflect.ValueOf(pkt).Elem().Type().Name())

	return pkt, nil
}

func (c *JavaClient) handle(pkt Packet) error {
	switch c.state {
	case Handshake:
		return c.handshakeHandlers(pkt)
	case Status:
		return c.statusHandlers(pkt)
	case Login:
		return c.loginHandlers(pkt)
	}
	return func(pkt Packet) error {
		panic("not implemented")
	}(nil)
}

func (c *JavaClient) disconnect(reason string) {
	if c.state == Login {
		// TODO: Send login kick player
	} else if c.state == Play {
		// TODO: Send play disconnect
	}
	_ = c.conn.Close()

	// TODO: Fix concurrent map writes (atomic or mutex this)
	delete(c.server.Clients, c.networkId)
}

func (c *JavaClient) disconnectError(err error) {
	// TODO: Better error kick handling
	log.Println(err)
	c.disconnect(err.Error())
}

func (c *JavaClient) enableCompression() error {
	err := c.Send(&SetCompression{Threshold: int32(config.Get.CompressionThreshold)})
	if err != nil {
		return err
	}
	c.compressionEnabled = true
	return nil
}

func (c *JavaClient) enableEncryption(sharedSecret []byte) (err error) {
	block, err := aes.NewCipher(sharedSecret)
	if err != nil {
		return
	}

	c.sharedSecret = sharedSecret
	c.encryptStream, err = auth.NewCFB8Encrypter(block, sharedSecret)
	if err != nil {
		return
	}

	c.decryptStream, err = auth.NewCFB8Decrypter(block, sharedSecret)
	if err != nil {
		return
	}

	c.reader = bufio.NewReader(cipher.StreamReader{
		S: c.decryptStream,
		R: c.conn,
	})

	c.writer = bufio.NewWriter(cipher.StreamWriter{
		S: c.encryptStream,
		W: c.conn,
	})

	c.encryptionEnabled = true
	return
}

func compress(data []byte) []byte {
	buf := &bytes.Buffer{}
	writer := zlib.NewWriter(buf)
	defer writer.Close()

	writer.Write(data)
	writer.Flush()
	return buf.Bytes()
}
