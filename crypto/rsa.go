package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func ParseRsaPrivKey(privKey string) (*rsa.PrivateKey, error) {
	var block *pem.Block
	decoded, err := hex.DecodeString(privKey)
	if err == nil {
		block, _ = pem.Decode(decoded)
	} else {
		block, _ = pem.Decode([]byte(privKey))
	}

	if block == nil {
		return nil, errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey.(*rsa.PrivateKey), nil
}

func ParseRsaPubKey(pubKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return nil, errors.New("public key error")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}
	rsaPub := pub.(*rsa.PublicKey)
	return rsaPub, nil
}

func RsaEncryptOAEP(pubKey *rsa.PublicKey, data []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, data, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

func RsaDecryptOAEP(privKey *rsa.PrivateKey, encryptedBytes []byte) ([]byte, error) {
	plainBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedBytes, []byte("HBC_MPC"))
	if err != nil {
		return nil, err
	}
	return plainBytes, nil
}
