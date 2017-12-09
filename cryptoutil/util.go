package cryptoutil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
)

func GetPrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err.Error())
	}
	return privateKey
}

func GetPEMPrivateKey(prKey *rsa.PrivateKey) []byte {
	prKeyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prKey),
	}

	return pem.EncodeToMemory(prKeyBlock)
}

func GetPEMPublicKey(puKey rsa.PublicKey) []byte {
	asn1Bytes, err := asn1.Marshal(puKey)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		return nil
	}
	var puKeyBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	return pem.EncodeToMemory(puKeyBlock)
}

func HashSha256(rawData []byte) []byte {
	sha_256 := sha256.New()
	sha_256.Write(rawData)
	return sha_256.Sum(nil)
}

func RSASign(priKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	// https://golang.org/pkg/crypto/rsa/#example_SignPKCS1v15
	rng := rand.Reader
	hashed := sha256.Sum256(data)
	signature, err := rsa.SignPKCS1v15(rng, priKey, crypto.SHA256, hashed[:])
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		return nil, err
	}
	return signature, nil
}

func RSAVerify(pubKey *rsa.PublicKey, data []byte, signature []byte) bool {
	// https://golang.org/pkg/crypto/rsa/#VerifyPKCS1v15
	// Only small messages can be signed directly; thus the hash of a
	// message, rather than the message itself, is signed.
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		fmt.Printf("Error from verification: %s\n", err)
		return false
	}
	return true
}
