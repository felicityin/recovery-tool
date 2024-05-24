package ckd

import (
	"encoding/hex"
	"math/big"
	"testing"

	edwards "github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/stretchr/testify/assert"

	"recovery-tool/crypto"
	"recovery-tool/crypto/common"
)

func TestDerivation(t *testing.T) {
	priKeyBytes, err := hex.DecodeString("ae1e5bf5f3d6bf58b5c222088671fcbe78b437e28fae944c793897b26091f241")
	assert.NoError(t, err)
	chainCode, err := hex.DecodeString("ae1e5bf5f3d6bf58b5c222088671fcbe78b437e28fae944c793897b26091f242")
	assert.NoError(t, err)
	priKey := new(big.Int).SetBytes(priKeyBytes)
	pubKey := crypto.ScalarBaseMult(edwards.Edwards(), priKey)

	childPrivKey, childPubKey, err := DeriveEddsaChildPrivKey(priKey, pubKey, pubKey, chainCode, "81/0/0/35/0")
	assert.NoError(t, err)

	childPubKeyPt, err := DeriveEddsaChildPubKey(pubKey, pubKey, chainCode, "81/0/0/35/0")
	assert.NoError(t, err)

	childPubKeyBytes := edwards.NewPublicKey(childPubKeyPt.X(), childPubKeyPt.Y()).Serialize()
	assert.Equal(t, childPubKey, childPubKeyBytes)

	sk, pk, err := edwards.PrivKeyFromScalar(common.PadToLengthBytesInPlace(childPrivKey[:], 32))
	assert.NoError(t, err)
	assert.Equal(t, pk.X, childPubKeyPt.X())

	data := make([]byte, 32)
	for i := range data {
		data[i] = byte(i)
	}

	r, s, err := edwards.Sign(sk, data)
	assert.NoError(t, err, "sign should not throw an error")

	pk1 := edwards.PublicKey{
		Curve: edwards.Edwards(),
		X:     childPubKeyPt.X(),
		Y:     childPubKeyPt.Y(),
	}
	ok := edwards.Verify(&pk1, data, r, s)
	assert.True(t, ok)
}
