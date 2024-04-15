package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParseRsaPrivKey(priKey string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(priKey))
	if block == nil {
		return nil, errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey.(*rsa.PrivateKey), nil
}

func RsaEncryptOAEP(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

func RSADecryptOAEP(privKey *rsa.PrivateKey, encryptedBytes []byte) ([]byte, error) {
	plainBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedBytes, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return plainBytes, nil
}
