package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/shopspring/decimal"

	"recovery-tool/common/code"
	"recovery-tool/tx/apt"
	"recovery-tool/tx/dot"
	"recovery-tool/tx/sol"
)

func Sign(chain, url, privkey, toAddr, amount, coinAddress string) (string, error) {
	priv, err := hex.DecodeString(privkey)
	if err != nil {
		return "", code.NewI18nError(code.PrivkeyIsHex, "The private key should be in hexadecimal format")
	}

	if toAddr == "" {
		return "", code.NewI18nError(code.DstAddrNotEmpty, "The target address cannot be empty, please re-enter.")
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
		sig, err := sol.Sign(coinAddress, priv, toAddr, amountDec)
		if err != nil {
			return "", err
		}
		return sig, nil
	case "apt":
		if url == SolNode || url == "" {
			url = AptNode
		}
		apt := apt.NewApt(url)
		sig, err := apt.Sign("", priv, toAddr, amountDec)
		if err != nil {
			return "", err
		}
		return sig, nil
	case "dot":
		if url == SolNode || url == "" {
			url = DotNode
		}
		dot := dot.NewDot(url)
		sig, err := dot.Sign("", priv, toAddr, amountDec)
		if err != nil {
			return "", err
		}
		return sig, nil
	default:
		return "", code.NewI18nError(code.ChainParamErr, fmt.Sprintf("Unsupported chain: %s", chain))
	}
}
