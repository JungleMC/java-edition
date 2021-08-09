package net

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

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
