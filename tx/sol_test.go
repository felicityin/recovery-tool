package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/blocto/solana-go-sdk/rpc"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"recovery-tool/tx/sol"
	so "recovery-tool/tx/sol/solana-adapter/common"
)

// faucet: https://faucet.quicknode.com/solana/testnet

func TestSolGetBalance(t *testing.T) {
	node := rpc.TestnetRPCEndpoint
	address := "HXx8Ky1aY7GBLUghbadKais5QHdeJfdQ7mmgR9j4sqNK"

	sol := sol.NewSol(node)

	balance, err := sol.GetBalance(context.Background(), address)
	fmt.Printf("sol balance: %v\n", balance)
	assert.NoError(t, err)

	associatedAddress, err := sol.GetAssociatedAddress(address, so.UsdtSolana)
	assert.NoError(t, err)

	decimals, amount, amountDecimal, err := sol.GetTokenBalance(associatedAddress)
	assert.NoError(t, err)

	fmt.Printf("token balance: %v\n", amountDecimal)
	fmt.Printf("token decimals: %d\n", decimals)
	fmt.Printf("token amount: %s\n", amount)

	assert.True(t, false)
}

func TestSolTransfer(t *testing.T) {
	node := rpc.TestnetRPCEndpoint
	sol := sol.NewSol(node)

	priv, _ := hex.DecodeString("078fe2333b309a95f8bc59f6e03a10c4b7b51f3e12b7ccd4a62c41363a08437a")
	toAddr := "DeQNVvKUsJpqX84YfTpEbd5EGMbR1MfcrkMD8zpyxG9K"
	amountDec, err := decimal.NewFromString("0.0001")
	if err != nil {
		t.Errorf("decimal.NewFromString err: %s", err.Error())
		return
	}

	txHash, err := sol.Transfer("", priv, toAddr, amountDec)
	assert.NoError(t, err)
	fmt.Printf("tx: %s", txHash)

	txHash, err = sol.Transfer(so.UsdtSolana, priv, toAddr, amountDec)
	assert.NoError(t, err)
	fmt.Printf("tx: %s", txHash)
}
