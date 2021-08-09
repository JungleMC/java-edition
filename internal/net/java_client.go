package net

import (
	"bytes"
	"crypto/cipher"
	"github.com/JungleMC/java-edition/internal/config"
	"io"
	"log"
	"net"
	"reflect"
)

const MTU = 1500

type JavaClient struct {
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
}

func (c *JavaClient) listen() {
	for {
		buf := make([]byte, MTU)
		bytesRead, err := c.connection.Read(buf)
		if err != nil && err != io.EOF {
			c.disconnectError(err)
			return
		}

		buf = buf[:bytesRead]
		if c.encryptionEnabled {
			c.decryptStream.XORKeyStream(buf, buf)
		}

		reader := bytes.NewBuffer(buf)
		pkt, err := readPacket(reader, c.protocol, c.compressionEnabled)
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

func (c *JavaClient) send(pkt Packet) error {
	// TODO: Submit packets to a FIFO queue before sending directly, maintaining packet order
	buf := &bytes.Buffer{}
	writePacket(buf, reflect.ValueOf(pkt).Elem(), c.protocol, c.compressionEnabled, config.Get.CompressionThreshold)

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
}

func (c *JavaClient) disconnectError(err error) {
	// TODO: Better error kick handling
	log.Println(err)
	c.disconnect(err.Error())
}
