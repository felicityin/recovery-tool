package sol

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
	"github.com/shopspring/decimal"

	cm "recovery-tool/common"
	"recovery-tool/common/code"
	"recovery-tool/tx/eddsa"
	"recovery-tool/tx/sol/solana-adapter/common"
)

func CalcAddress(privkey []byte) string {
	publicKey := eddsa.Pubkey(privkey)
	return base58.Encode(publicKey.Serialize())
}

type Sol struct {
	Client   *Client
	Decimals int
}

func NewSol(url string) *Sol {
	client := NewClient(url)

	sol := new(Sol)
	sol.Client = client
	sol.Decimals = 9
	return sol
}

func (c *Sol) GetAssociatedAddress(ownerAddres string, coinAddress string) (associatedAddress string, err error) {
	if coinAddress == "" {
		return
	}

	var res GetTokenAccountsByOwnerResponse
	res, err = c.Client.GetTokenAccountsByOwnerWithCfg(context.Background(), ownerAddres, coinAddress, Cfg{Encoding: "jsonParsed"})
	if err != nil {
		cm.Logger.Errorf("GetTokenAccountsByOwnerWithCfg err: %s", err.Error())
		return
	}

	contractBalanceDecimal := decimal.Zero
	for _, v := range res.Result.Value {
		if v.Account.Data.Program == "spl-token" {
			amountDecimal, _ := decimal.NewFromString(v.Account.Data.Parsed.Info.TokenAmount.UIAmountString)
			if amountDecimal.GreaterThan(contractBalanceDecimal) {
				contractBalanceDecimal = amountDecimal
				associatedAddress = v.Pubkey
			}
		}
	}

	if contractBalanceDecimal.IsZero() && len(res.Result.Value) > 0 {
		associatedAddress = res.Result.Value[0].Pubkey
	}
	return
}

func (c *Sol) GetContractDecimals(coinAddress string) (decimals int64, err error) {
	res, err := c.Client.GetTokenSupply(context.Background(), coinAddress)
	if err != nil {
		err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
		return 0, err
	}
	return int64(res.Result.Value.Decimals), nil
}

func (c *Sol) GetBalance(ctx context.Context, base58Addr string) (uint64, error) {
	return c.Client.GetBalance(ctx, base58Addr)
}

func (c *Sol) GetTokenBalance(associatedAddress string) (decimals int, amount string, amountDecimal decimal.Decimal, err error) {
	res, err1 := c.Client.GetAccountInfoWithCfg(context.Background(), associatedAddress, Cfg{
		Encoding: "jsonParsed",
	})
	if err1 != nil {
		d, _ := decimal.NewFromString("0")
		if strings.Contains(err1.Error(), "WrongSize") {
			return 0, "0", d, nil
		}
		cm.Logger.Errorf("GetAccountInfoWithCfg: %s", err1.Error())
		err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
		return
	}

	var data = struct {
		Parsed struct {
			Info struct {
				IsNative    bool   `json:"isNative"`
				Mint        string `json:"mint"`
				Owner       string `json:"owner"`
				State       string `json:"state"`
				TokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UIAmount       float64 `json:"uiAmount"`
					UIAmountString string  `json:"uiAmountString"`
				} `json:"tokenAmount"`
			} `json:"info"`
			Type string `json:"type"`
		} `json:"parsed"`
		Program string `json:"program"`
		Space   int    `json:"space"`
	}{}

	err = json.Unmarshal(res.Result.Value.Data, &data)
	if err != nil {
		err = fmt.Errorf("unmarshal token balance res err: %s", err.Error())
		return
	}

	decimals = data.Parsed.Info.TokenAmount.Decimals
	amount = data.Parsed.Info.TokenAmount.Amount
	amountDecimal, _ = decimal.NewFromString(data.Parsed.Info.TokenAmount.UIAmountString)
	return
}

func (c *Sol) Sign(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (sig string, err error) {
	blockHash, err := c.GetBlockHash()
	if err != nil {
		cm.Logger.Errorf("get block hash err: %s", err.Error())
		err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
		return
	}

	fromAddr := CalcAddress(privkey)
	cm.Logger.Infof("from: %s", fromAddr)

	properties := map[string]interface{}{
		"from":             fromAddr,
		"to":               toAddr,
		"decimals":         "9",
		"token_program_id": "",
		"type":             "",
		"gas_price":        "500",
		"gas_limit":        "400000",
		"block_hash":       blockHash,
	}

	var decimals int64

	if coinAddress != "" {
		if coinAddress == common.CwifSolana {
			properties["token_program_id"] = "TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb"
		}

		fromAssociatedAddress, err1 := c.GetAssociatedAddress(fromAddr, coinAddress)
		if err1 != nil {
			cm.Logger.Errorf("get from associated addr err: %s", err1.Error())
			if strings.Contains(err.Error(), "Invalid param: WrongSize") {
				return "", code.NewI18nError(code.CoinAddrInvalid, "The token address format is wrong, please re-enter.")
			}
			err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
			return
		}
		if fromAssociatedAddress == "" {
			err = code.NewI18nError(code.SrcCoinAccountNotFound, "The sending token address does not exist, please check and try again.")
			return
		}
		properties["from"] = fromAssociatedAddress
		cm.Logger.Infof("from associated addr: %s", fromAssociatedAddress)

		toAssociatedAddress, err1 := c.GetAssociatedAddress(toAddr, coinAddress)
		if err1 != nil {
			cm.Logger.Errorf("get to associated addr err: %s", err1.Error())
			err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
			return
		}
		cm.Logger.Infof("to associated addr: %s", toAssociatedAddress)
		if toAssociatedAddress == "" {
			properties["type"] = SignTypeAccount
		} else {
			properties["to"] = toAssociatedAddress
		}

		decimals, err = c.GetContractDecimals(coinAddress)
		if err != nil {
			cm.Logger.Errorf("get decimals err: %s", err1.Error())
			err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
			return
		}
		properties["decimals"] = fmt.Sprintf("%d", decimals)
	}

	solPacker := new(solanaPackager)

	err = solPacker.Pack("sol", properties, coinAddress, amountDec, fromAddr, toAddr)
	if err != nil {
		return
	}

	txRaw, err := solPacker.GetRaw()
	if err != nil {
		err = fmt.Errorf("GetRaw err: %s", err.Error())
		return
	}

	msgBytes, err := txRaw.Serialize()
	if err != nil {
		err = fmt.Errorf("GetSignPacket err: %s", err.Error())
		return
	}

	signature, err := eddsa.Sign(privkey, msgBytes)
	if err != nil {
		return
	}

	tx := types.Transaction{
		Signatures: []types.Signature{signature},
		Message:    *txRaw,
	}
	serialize, err := tx.Serialize()
	if err != nil {
		return
	}

	return base58.Encode(serialize), nil
}

func (c *Sol) Transfer(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (txHash string, err error) {
	sig, err := c.Sign(coinAddress, privkey, toAddr, amountDec)
	if err != nil {
		return
	}

	response, err := c.Client.SendTransaction(context.Background(), sig)
	if err != nil {
		cm.Logger.Errorf("send tx err: %s", err.Error())
		return
	}
	return response.Result, nil
}

func (c *Sol) GetBlockHash() (blockHash string, err error) {
	var res GetLatestBlockHashResponse
	var res2 GetBlockHeightResponse
	for i := 0; i < 5; i++ {
		res, err = c.Client.GetLatestBlockHash(context.Background())
		if err != nil {
			cm.Logger.Errorf("get block hash err: %s", err.Error())
			err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
			return
		}
		lastValidBlockHeight := res.Result.Value.LastValidBlockHeight + 150
		res2, err = c.Client.GetBlockHeight(context.Background())
		if err != nil {
			cm.Logger.Errorf("get block height err: %s", err.Error())
			err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
			return
		}
		blockHeight := res2.Result
		if blockHeight < lastValidBlockHeight {
			blockHash = res.Result.Value.Blockhash
			return
		}
		time.Sleep(time.Second * 10) // wait 10s
	}
	return res.Result.Value.Blockhash, nil
}
