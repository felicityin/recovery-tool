package eddsa

import (
	"crypto/sha512"
	"fmt"
	"math/big"

	e "github.com/agl/ed25519/edwards25519"
	"github.com/decred/dcrd/dcrec/edwards/v2"

	"recovery-tool/common"
	"recovery-tool/crypto"
)

func Pubkey(privkey []byte) *edwards.PublicKey {
	pubECPoint := crypto.ScalarBaseMult(edwards.Edwards(), new(big.Int).SetBytes(privkey))
	return edwards.NewPublicKey(pubECPoint.X(), pubECPoint.Y())
}

func Sign(privkey []byte, message []byte) ([]byte, error) {
	seed := common.GetRandomPositiveInt(edwards.Edwards().Params().N)
	h := sha512.New()
	h.Reset()
	h.Write(seed.Bytes()[:])
	h.Write(privkey[:])
	h.Write(message)

	var messageDigest [64]byte
	h.Sum(messageDigest[:0])
	var messageDigestReduced [32]byte
	e.ScReduce(&messageDigestReduced, &messageDigest)

	var R e.ExtendedGroupElement
	e.GeScalarMultBase(&R, &messageDigestReduced)

	// compute lambda
	var encodedR [32]byte
	R.ToBytes(&encodedR)

	pubkey := Pubkey(privkey)
	encodedPubKey := EcPointToEncodedBytes(pubkey.X, pubkey.Y)

	// h = hash512(k || A || M)
	var hramDigest [64]byte
	h.Reset()
	h.Write(encodedR[:])
	h.Write(encodedPubKey[:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	e.ScReduce(&hramDigestReduced, &hramDigest)

	// compute s
	var s [32]byte
	e.ScMulAdd(&s, &hramDigestReduced, BigIntToEncodedBytes(new(big.Int).SetBytes(privkey)), &messageDigestReduced)

	r := EncodedBytesToBigInt(&encodedR)
	si := EncodedBytesToBigInt(&s)

	if !edwards.Verify(pubkey, message, r, si) {
		return nil, fmt.Errorf("edwards.Verify failed")
	}

	sig := append(encodedR[:], s[:]...)
	return sig, nil
}
