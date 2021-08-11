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
	RDB        *redis.Client
	privateKey *rsa.PrivateKey
	privateKeyBytes []byte
	publicKeyBytes  []byte

	Clients map[uuid.UUID]*JavaClient
}

func (s *JavaServer) GetClient(networkId []byte) (*JavaClient, error) {
	id, err := uuid.FromBytes(networkId)
	if err != nil {
		return nil, err
	}
	return s.Clients[id], nil
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

	s.Clients[networkId] = client

	client.listen()
}

func (s *JavaServer) Bootstrap(rdb *redis.Client, address string, port int, onlineMode bool) error {
	if onlineMode {
		s.generateKeyPair()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
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
