package net

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

func (s *JavaServer) generateKeyPair() {
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}

	privKey.Precompute()
	if err := privKey.Validate(); err != nil {
		panic(err)
	}
	s.privateKey = privKey
	s.privateKeyBytes = x509.MarshalPKCS1PrivateKey(privKey)
	s.publicKeyBytes, _ = x509.MarshalPKIXPublicKey(privKey.Public())
}
