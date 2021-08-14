package net

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net/auth"
	"github.com/JungleMC/java-edition/internal/net/packets"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"reflect"
)

const MTU = 1500

type JavaClient struct {
	networkId           uuid.UUID
	server              *JavaServer
	connection          net.Conn
	protocol            Protocol
	gameProtocolVersion int32
	compressionEnabled  bool
	encryptionEnabled   bool
	verifyToken         []byte
	sharedSecret        []byte
	encryptStream       cipher.Stream
	decryptStream       cipher.Stream
	authProfile         *auth.Profile
}

func (c *JavaClient) listen() {
	for {
		buf := make([]byte, MTU)
		bytesRead, err := c.connection.Read(buf)
		if err != nil && err != io.EOF {
			c.disconnectError(err)
			return
		}

		if bytesRead == 0 { // EOF
			continue
		}

		buf = buf[:bytesRead]
		if c.encryptionEnabled {
			c.decryptStream.XORKeyStream(buf, buf)
		}

		reader := bytes.NewBuffer(buf)
		pkt, err := ReadPacket(reader, c.protocol, c.compressionEnabled)
		if err != nil && err != io.EOF {
			c.disconnectError(err)
			return
		}

		if err == io.EOF {
			c.disconnect("")
			return
		}

		err = c.handle(pkt)
		if err != nil {
			c.disconnectError(err)
		}
	}
}

func (c *JavaClient) Send(pkt Packet) error {
	if pkt == nil {
		return nil
	}

	if config.Get.Verbose {
		log.Printf("tx -> %v\n", reflect.ValueOf(pkt).Elem().Type().Name())
	}

	// TODO: Submit packets to a FIFO queue before sending directly, maintaining packet order
	buf := &bytes.Buffer{}
	WritePacket(buf, reflect.ValueOf(pkt).Elem(), c.protocol, c.compressionEnabled, config.Get.CompressionThreshold)

	data := buf.Bytes()
	if config.Get.OnlineMode && c.encryptionEnabled {
		c.encryptStream.XORKeyStream(data, data)
	}

	_, err := c.connection.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *JavaClient) handle(pkt Packet) error {
	if pkt == nil {
		return nil
	}

	if config.Get.Verbose {
		log.Printf("rx <- %v\n", reflect.ValueOf(pkt).Type().Name())
	}

	switch c.protocol {
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
	if c.protocol == Login {
		// TODO: Send login kick player
	} else if c.protocol == Play {
		// TODO: Send play disconnect
	}
	_ = c.connection.Close()

	// TODO: Fix concurrent map writes (atomic or mutex this)
	delete(c.server.Clients, c.networkId)
}

func (c *JavaClient) disconnectError(err error) {
	// TODO: Better error kick handling
	log.Println(err)
	c.disconnect(err.Error())
}

func (c *JavaClient) enableCompression() error {
	err := c.Send(&packets.ClientboundLoginCompressionPacket{Threshold: int32(config.Get.CompressionThreshold)})
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

	c.encryptionEnabled = true
	return
}
