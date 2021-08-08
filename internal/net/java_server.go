package net

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net"
)

type JavaServer struct {
	RDB             *redis.Client
	privateKey      *rsa.PrivateKey
	privateKeyBytes []byte
	publicKeyBytes  []byte
}

func (s *JavaServer) clientConnect(connection net.Conn) {
	client := &JavaClient{
		server: s,
		connection: connection,
		protocol:   Handshake,
	}
	_, _ = rand.Read(client.verifyToken)

	client.listen()
}

func Bootstrap(rdb *redis.Client, address string, port int, onlineMode bool) (*JavaServer, error) {
	s := &JavaServer{
		RDB: rdb,
	}

	if onlineMode {
		s.generateKeyPair()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	for {
		connection, clientErr := listener.Accept()

		if clientErr != nil {
			log.Println(err)
			connection.Close()
			continue
		}

		go s.clientConnect(connection)
	}
}
