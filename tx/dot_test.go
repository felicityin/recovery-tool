package tx

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"recovery-tool/tx/dot"
)

func TestDotTransfer(t *testing.T) {
	node := "https://polkadot-public-sidecar.parity-chains.parity.io"
	dot := dot.NewDot(node)

	priv, _ := hex.DecodeString("0d9e11aeaa5d1f00565386799fa6e04e51c3b8087113e972d0dfc4bcc26ad9dc")
	toAddr := "13fHvYjRAC1Ebj1LpfiaBC5jVgS2wvVxXziPnGEzCNJPwvBk"
	amountDec, err := decimal.NewFromString("0.0001")
	if err != nil {
		t.Errorf("decimal.NewFromString err: %s", err.Error())
		return
	}

	txHash, err := dot.Transfer("", priv, toAddr, amountDec)
	assert.NoError(t, err)
	fmt.Printf("tx: %s", txHash)
}
