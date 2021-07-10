package net

import (
	"bytes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"net"
)

const MTU = 1500

type JavaClient struct {
	connection         net.Conn
	protocol           Protocol
	compressionEnabled bool
	encryptionEnabled  bool
	verifyToken        []byte
	sharedSecret       []byte
	encryptStream      cipher.Stream
	decryptStream      cipher.Stream
}

func clientConnect(connection net.Conn) {
	client := &JavaClient{
		connection: connection,
		protocol:   Handshake,
	}
	_, _ = rand.Read(client.verifyToken)

	client.listen()
}

func (c *JavaClient) listen() {
	for {
		buf := make([]byte, MTU)
		bytesRead, err := c.connection.Read(buf)
		if err != nil && err != io.EOF {
			c.disconnect()
			return
		}

		buf = buf[:bytesRead]
		if c.encryptionEnabled {
			c.decryptStream.XORKeyStream(buf, buf)
		}

		reader := bytes.NewBuffer(buf)
		pkt, err := readPacket(reader, c.protocol, c.compressionEnabled)
		if err != nil && err != io.EOF {
			c.disconnect()
			return
		}

		if pkt == nil {
			continue
		}

		c.receivePacket(pkt)
	}
}

func (c *JavaClient) receivePacket(pkt Packet) {
}

func (c *JavaClient) disconnect() {
	if c.protocol == Play {
	}
	_ = c.connection.Close()
}
