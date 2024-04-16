package common_test

import (
	"encoding/hex"
	"recovery-tool/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMnemonic(t *testing.T) {
	mnemonicWords := "amused garlic window please enrich sick gate ready owner giraffe elite umbrella hair seat punch seminar notable enroll wet asset outdoor inflict rich mushroom"
	priv, _, err := common.CalcMasterPriv(mnemonicWords)
	privKey := hex.EncodeToString(priv[:])
	assert.NoError(t, err)
	assert.Equal(t, privKey, "9aa5eaa7c63f1e157b94896919dd4327d279b8265c4794100f28b78cae79be7d")
}
