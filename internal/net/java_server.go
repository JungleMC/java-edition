package net

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"github.com/JungleMC/java-edition/internal/net/auth"
	. "github.com/JungleMC/protocol"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"net"
)

type JavaServer struct {
	RDB             *redis.Client
	privateKey      *rsa.PrivateKey
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

func (s *JavaServer) clientConnect(conn net.Conn) {
	networkId, _ := uuid.NewRandom()

	client := &JavaClient{
		networkId: networkId,
		server:    s,
		conn:      conn,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
		state:     Handshake,
		authProfile: &auth.Profile{
			ID: uuid.New(), // Offline mode. Online mode overrides this
		},
		verifyToken: make([]byte, 4),
	}
	_, _ = rand.Read(client.verifyToken)

	s.Clients[networkId] = client
	client.listen()
}

func (s *JavaServer) generateKeyPair() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	privateKey.Precompute()
	if err = privateKey.Validate(); err != nil {
		panic(err)
	}
	s.privateKey = privateKey
	s.privateKeyBytes = x509.MarshalPKCS1PrivateKey(privateKey)
	s.publicKeyBytes, _ = x509.MarshalPKIXPublicKey(privateKey.Public())
}

func (s *JavaServer) Bootstrap(address string, port int, onlineMode bool) error {
	if onlineMode {
		s.generateKeyPair()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, clientErr := listener.Accept()
		if clientErr != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		go s.clientConnect(conn)
	}
}
