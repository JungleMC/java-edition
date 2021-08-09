package net

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/JungleMC/java-edition/internal/net/auth"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"net"
)

type JavaServer struct {
	rdb        *redis.Client
	privateKey *rsa.PrivateKey
	privateKeyBytes []byte
	publicKeyBytes  []byte

	clients map[uuid.UUID]*JavaClient
}

func (s *JavaServer) clientConnect(connection net.Conn) {
	networkId, _ := uuid.NewRandom()

	client := &JavaClient{
		networkId:  networkId,
		server:     s,
		connection: connection,
		protocol:   Handshake,
		authProfile: &auth.Profile{},
		verifyToken: make([]byte, 4),
	}
	_, _ = rand.Read(client.verifyToken)

	s.clients[networkId] = client

	client.listen()
}

func Bootstrap(rdb *redis.Client, address string, port int, onlineMode bool) (*JavaServer, error) {
	s := &JavaServer{
		rdb: rdb,
		clients: make(map[uuid.UUID]*JavaClient),
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
