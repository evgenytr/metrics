package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/evgenytr/metrics.git/internal/config"
)

func main() {

	publicKeyFileStoragePath, privateKeyFileStoragePath, err := config.GetRsaKeyGeneratorConfig()
	fmt.Println(publicKeyFileStoragePath, privateKeyFileStoragePath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := &privateKey.PublicKey

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if err != nil {
		log.Fatal(err)
	}

	var publicKeyPEM bytes.Buffer
	err = pem.Encode(&publicKeyPEM, &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(publicKey),
	})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(privateKeyFileStoragePath)
	err = os.WriteFile(privateKeyFileStoragePath, privateKeyPEM.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(publicKeyFileStoragePath, publicKeyPEM.Bytes(), 0666)
	if err != nil {
		log.Fatal(err)
	}

}
