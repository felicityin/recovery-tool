package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"

	"recovery-tool/common"
	"recovery-tool/common/code"
	"recovery-tool/tx/apt"
	"recovery-tool/tx/dot"
	"recovery-tool/tx/sol"
	"recovery-tool/tx/ton"
)

func Transfer(chain, url, privkey, toAddr, amount, coinAddress string) (string, error) {
	priv, err := hex.DecodeString(privkey)
	if err != nil {
		return "", code.NewI18nError(code.PrivkeyInvalid, "The private key format is wrong, please re-enter.")
	}

	if toAddr == "" {
		return "", code.NewI18nError(code.DstAddrNotEmpty, "The target address cannot be empty, please re-enter.")
	}

	amountDec, err := decimal.NewFromString(amount)
	if err != nil {
		return "", code.NewI18nError(code.AmountInvalid, "Unable to convert transfer amount to decimal")
	}

	if chain != "sol" && chain != "ton" && coinAddress != "" {
		common.Logger.Errorf("chain: %s", chain)
		return "", code.NewI18nError(code.CoinUnsupported, "This chain only supports main chain coins for now")
	}

	switch chain {
	case "sol":
		sol := sol.NewSol(url)
		txHash, err := sol.Transfer(coinAddress, priv, toAddr, amountDec)
		if err != nil {
			common.Logger.Errorf("[sol] transfer err: %s", err.Error())
			if strings.Contains(err.Error(), "insufficient lamports") || strings.Contains(err.Error(), "account (0) with insufficient funds for rent") {
				return "", code.NewI18nError(code.SolInsufficientFunds, "Insufficient gas fee (the current maximum transaction fee on the chain is 0.00089608 sol).")
			}
			if strings.Contains(err.Error(), "AccountNotFound") {
				return "", code.NewI18nError(code.SrcAccountNotFound, "The sending account does not exist, please check and try again")
			}
			if strings.Contains(err.Error(), "account (1) with insufficient funds for rent") {
				return "", code.NewI18nError(code.DstAccountNotFound, "The receiving account does not exist, please check and try again")
			}
			if strings.Contains(err.Error(), "get status code: 429") || strings.Contains(err.Error(), "Blockhash not found") {
				return "", code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
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
			if strings.Contains(err.Error(), "INSUFFICIENT_BALANCE_FOR_TRANSACTION_FEE") {
				return "", code.NewI18nError(code.AptInsufficientFunds, "Insufficient gas fee (the current maximum transaction fee on the chain is 0.002 apt).")
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
			if strings.Contains(err.Error(), "Inability to pay some fees") {
				return "", code.NewI18nError(code.DotInsufficientFunds, "Insufficient gas fee (the current maximum transaction fee on the chain is 1 dot).")
			}
			return "", err
		}
		return txHash, nil
	case "ton":
		if url == SolNode || url == "" {
			url = TonNode
		}
		ton := ton.NewTon(url)
		txHash, err := ton.Transfer(coinAddress, priv, toAddr, amountDec)
		if err != nil {
			common.Logger.Errorf("[ton] transfer err: %s", err.Error())
			if strings.Contains(err.Error(), "get status code: 429") {
				return "", code.NewI18nError(code.TonNetworkErr, "Using API without API key is limited to 1 request per second. Register your API key in the https://toncenter.com to get access with higher limits.")
			}
			if strings.Contains(err.Error(), "failed to run get_wallet_address method") {
				return "", code.NewI18nError(code.CoinAddrNotExists, "The counterparty contract address does not exist, please re-enter it.")
			}
			return "", err
		}
		fmt.Printf("tx hash: %s\n", txHash)
		txId, err := base64.StdEncoding.DecodeString(txHash)
		if err != nil {
			return "", fmt.Errorf("base64 decode tx hash err: %s", err.Error())
		}
		return hex.EncodeToString(txId), nil
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
	case "ton":
		return TonScan
	default:
		return ""
	}
}
