package eddsa

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	msg := []byte("0f1bce53a4ec6ab3d62726aeb51b10bc534878200a83633bba084a9cb63fc1080f1bce53a4ec6ab3d62726aeb51b10bc534878200a83633bba084a9cb63fc1080f1bce53a4ec6ab3d62726aeb51b10bc534878200a83633bba084a9cb63fc108")
	priv, _ := hex.DecodeString("078fe2333b309a95f8bc59f6e03a10c4b7b51f3e12b7ccd4a62c41363a08437a")
	_, err := Sign(priv, msg)
	assert.NoError(t, err)
}
