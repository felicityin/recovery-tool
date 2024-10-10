package tx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"recovery-tool/tx/apt"
)

func TestAptTransfer(t *testing.T) {
	node := "https://fullnode.devnet.aptoslabs.com"
	apt := apt.NewApt(node)

	priv, _ := hex.DecodeString("0f2020d5f3ff4a08919d6e5f9058c47946dffaa620ea10e0d884c078dfa6ba23")
	toAddr := "9b71cd285c0233d4f6db4c9c4ec3359cd0ab4e892b493e1416b52572bbdef8c8"
	amountDec, err := decimal.NewFromString("0.0001")
	if err != nil {
		t.Errorf("decimal.NewFromString err: %s", err.Error())
		return
	}

	txHash, err := apt.Transfer("", priv, toAddr, amountDec)
	assert.NoError(t, err)
	fmt.Printf("tx: %s", txHash)
}
