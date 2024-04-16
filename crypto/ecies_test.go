package crypto

import (
	"testing"

	ecies "github.com/ecies/go/v2"
	"github.com/stretchr/testify/assert"
)

func TestEcies(t *testing.T) {
	privkey := "ea5db436b7508e5c8ec3ae17003bcb997c30e03c655f0dd2d1824ec93bd0501c"
	pubkey := "0232dbc41db7fc649ca5aac4456566cb750bba3f025d8ef2fd93eec60bb07f85e0"

	eciesPrivkey, err := ecies.NewPrivateKeyFromHex(privkey)
	assert.NoError(t, err)

	eciesPubkey, err := ecies.NewPublicKeyFromHex(pubkey)
	assert.NoError(t, err)

	assert.Equal(t, eciesPrivkey.PublicKey, eciesPubkey)

	encrypted, err := ecies.Encrypt(eciesPubkey, []byte("hello"))
	assert.NoError(t, err)

	decrypted, err := ecies.Decrypt(eciesPrivkey, encrypted)
	assert.NoError(t, err)

	assert.Equal(t, decrypted, []byte("hello"))
}
