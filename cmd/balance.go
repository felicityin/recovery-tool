package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/shopspring/decimal"

	"recovery-tool/common"
	"recovery-tool/common/code"
	"recovery-tool/tx/apt"
	"recovery-tool/tx/dot"
	"recovery-tool/tx/sol"
	"recovery-tool/tx/ton"
)

const (
	SolScan = "https://solscan.io/tx"
	AptScan = "https://aptoscan.com/transaction"
	DotScan = "https://polkadot.subscan.io/extrinsic"
	TonScan = "https://tonviewer.com/transaction"

	SolNode = "https://api.mainnet-beta.solana.com"
	AptNode = "https://fullnode.mainnet.aptoslabs.com"
	DotNode = "https://polkadot-asset-hub-public-sidecar.parity-chains.parity.io"
	TonNode = "https://toncenter.com"
)

type BalanceResult struct {
	Balance  string `json:"balance"`
	Decimals string `json:"decimals"`
	Amount   string `json:"amount"`
}

func (result BalanceResult) ToJsonStr() string {
	b, _ := json.Marshal(result)
	return string(b)
}

func GetBalance(chain, url, addr, coinAddress string) (*BalanceResult, error) {
	if addr == "" {
		return nil, code.NewI18nError(code.SrcAddrNotEmpty, "The address cannot be empty")
	}

	if chain != "sol" && chain != "ton" && coinAddress != "" {
		return nil, code.NewI18nError(code.CoinUnsupported, "This chain only supports main chain coins for now")
	}

	switch chain {
	case "sol":
		sol := sol.NewSol(url)
		if coinAddress == "" {
			balance, err := sol.GetBalance(context.Background(), addr)
			if err != nil {
				return nil, fmt.Errorf("get balance err: %s", err.Error())
			}
			amount, err := decimal.NewFromString(fmt.Sprintf("%d", balance))
			if err != nil {
				return nil, fmt.Errorf("amount to decimal err: %s", err.Error())
			}
			return &BalanceResult{
				Balance:  fmt.Sprintf("%d", balance),
				Decimals: "9",
				Amount:   fmt.Sprintf("%v", amount.Div(decimal.NewFromInt32(1000000000))),
			}, nil
		} else {
			associatedAddress, err := sol.GetAssociatedAddress(addr, coinAddress)
			if err != nil {
				return nil, fmt.Errorf("get associated address err: %s", err.Error())
			}
			decimals, balance, amount, err := sol.GetTokenBalance(associatedAddress)
			if err != nil {
				return nil, fmt.Errorf("get balance err: %s", err.Error())
			}
			return &BalanceResult{
				Balance:  balance,
				Decimals: fmt.Sprintf("%d", decimals),
				Amount:   fmt.Sprintf("%v", amount),
			}, nil
		}
	case "apt":
		if url == SolNode || url == "" {
			url = AptNode
		}
		apt := apt.NewApt(url)
		balance, amount, err := apt.Balance(addr)
		if err != nil {
			return nil, fmt.Errorf("get balance err: %s", err.Error())
		}
		return &BalanceResult{
			Balance:  fmt.Sprintf("%v", balance),
			Decimals: "8",
			Amount:   amount,
		}, nil
	case "dot":
		if url == SolNode || url == "" {
			url = DotNode
		}
		dot := dot.NewDot(url)
		balance, amount, err := dot.Balance(addr)
		if err != nil {
			return nil, fmt.Errorf("get balance err: %s", err.Error())
		}
		return &BalanceResult{
			Balance:  fmt.Sprintf("%v", balance),
			Decimals: "10",
			Amount:   amount,
		}, nil
	case "ton":
		if url == SolNode || url == "" {
			url = TonNode
		}
		ton, err := ton.NewTon(url)
		if err != nil {
			common.Logger.Errorf("[ton] create ton client err: %s", err.Error())
			return nil, err
		}
		if coinAddress == "" {
			balance, amount, err := ton.Balance(addr)
			if err != nil {
				return nil, fmt.Errorf("get balance err: %s", err.Error())
			}
			return &BalanceResult{
				Balance:  balance,
				Decimals: "9",
				Amount:   fmt.Sprintf("%v", amount),
			}, nil
		} else {
			ton.ContractAddress = coinAddress
			decimals, balance, amount, err := ton.BalanceOf(addr)
			if err != nil {
				return nil, fmt.Errorf("get balance err: %s", err.Error())
			}
			return &BalanceResult{
				Balance:  balance,
				Decimals: fmt.Sprintf("%d", decimals),
				Amount:   fmt.Sprintf("%v", amount),
			}, nil
		}
	default:
		return nil, code.NewI18nError(code.ChainParamErr, fmt.Sprintf("Unsupported chain: %s", chain))
	}
}

func ShortChainName(chain string) string {
	switch chain {
	case common.SolanaChain:
		return "sol"
	case common.AptostChain:
		return "apt"
	case common.PolkadotChain:
		return "dot"
	case common.TonChain:
		return "ton"
	default:
		return ""
	}
}
