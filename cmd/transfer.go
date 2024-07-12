package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"recovery-tool/common"
	"recovery-tool/common/code"
	"recovery-tool/tx/apt"
	"recovery-tool/tx/dot"
	"recovery-tool/tx/sol"
)

func Transfer(chain, url, privkey, toAddr, amount, coinAddress string) (string, error) {
	priv, err := hex.DecodeString(privkey)
	if err != nil {
		return "", code.NewI18nError(code.PrivkeyInvalid, "The private key should be in hexadecimal format")
	}

	if toAddr == "" {
		return "", code.NewI18nError(code.DstAddrNotEmpty, "The recipient's address cannot be empty")
	}

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return "", code.NewI18nError(code.AmountInvalid, "Unable to convert transfer amount to decimal")
	}

	if chain != "sol" && coinAddress != "" {
		return "", code.NewI18nError(code.CoinUnsupported, "This chain only supports main chain coins for now")
	}

	switch chain {
	case "sol":
		sol := sol.NewSol(url)
		txHash, err := sol.Transfer(coinAddress, priv, toAddr, amountDec)
		if err != nil {
			common.Logger.Errorf("[sol] transfer err: %s", err.Error())
			if strings.Contains(err.Error(), "insufficient") {
				return "", code.NewI18nError(code.SolInsufficientFunds, "Insufficient balance to pay for transaction fee. The max tx fee is 0.00089608 sol")
			}
			return "", err
		}
		return txHash, nil
	case "apt":
		if url == SolNode || url == "" {
			url = AptNode
		}
		apt := apt.NewApt(url)
		txHash, err := apt.Transfer("", priv, toAddr, amountDec)
		if err != nil {
			common.Logger.Errorf("[apt] transfer err: %s", err.Error())
			if strings.Contains(err.Error(), "INSUFFICIENT") {
				return "", code.NewI18nError(code.AptInsufficientFunds, "Insufficient balance to pay for transaction fee. The max tx fee is 0.002 apt")
			}
			return "", err
		}
		return txHash, nil
	case "dot":
		if url == SolNode || url == "" {
			url = DotNode
		}
		dot := dot.NewDot(url)
		txHash, err := dot.Transfer("", priv, toAddr, amountDec)
		if err != nil {
			common.Logger.Errorf("[dot] transfer err: %s", err.Error())
			if strings.Contains(err.Error(), "low") {
				return "", code.NewI18nError(code.DotInsufficientFunds, "Insufficient balance to pay for transaction fee. The max tx fee is 1 dot")
			}
			return "", err
		}
		return txHash, nil
	default:
		return "", code.NewI18nError(code.ChainParamErr, fmt.Sprintf("Unsupported chain: %s", chain))
	}
}

func Scan(chain string) string {
	switch chain {
	case "sol":
		return SolScan
	case "apt":
		return AptScan
	case "dot":
		return DotScan
	default:
		return ""
	}
}
