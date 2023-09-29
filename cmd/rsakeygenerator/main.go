package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {

	privateKeyFileStoragePath := "rsakeys/private.pem"
	publicKeyFileStoragePath := "rsakeys/public.pub"

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey

	var privateKeyPEM bytes.Buffer
	pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	var publicKeyPEM bytes.Buffer
	pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})

	err = os.WriteFile(privateKeyFileStoragePath, privateKeyPEM.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(publicKeyFileStoragePath, publicKeyPEM.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

}
