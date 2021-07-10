package net

import (
	"crypto/rsa"
	"fmt"
	"log"
	"net"
)

type JavaServer struct {
	privateKey      *rsa.PrivateKey
	privateKeyBytes []byte
	publicKeyBytes  []byte
}

func Bootstrap(address string, port int, onlineMode bool) (*JavaServer, error) {
	s := &JavaServer{}

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
		go clientConnect(connection)
	}
}
