package apt

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/decred/dcrd/dcrec/edwards/v2"
	"github.com/portto/aptos-go-sdk/client"
	crypt "github.com/portto/aptos-go-sdk/crypto"
	"github.com/portto/aptos-go-sdk/models"
	"github.com/shopspring/decimal"
	"github.com/the729/lcs"

	"recovery-tool/common"
	"recovery-tool/common/code"
	"recovery-tool/crypto"
	"recovery-tool/tx/eddsa"
)

func CalcAddress(privkey []byte) string {
	pubECPoint := crypto.ScalarBaseMult(edwards.Edwards(), new(big.Int).SetBytes(privkey))
	publicKey := edwards.NewPublicKey(pubECPoint.X(), pubECPoint.Y())
	accountAddress := models.AccountAddress{}
	accountAddress = crypt.SingleSignerAuthKey(publicKey.Serialize())
	return accountAddress.PrefixZeroTrimmedHex()
}

type Apt struct {
	Client          client.AptosClient
	GasLimit        decimal.Decimal
	Decimals        int
	ContractAddress string
	MaxFee          string
}

func NewApt(url string) *Apt {
	apt := new(Apt)
	apt.Client = client.NewAptosClient(url)
	apt.GasLimit = decimal.NewFromFloat(2000)
	apt.Decimals = 8
	apt.MaxFee = "0.02"
	return apt
}

func (c *Apt) GetBlockNumber() (uint64, error) {
	res, err := c.Client.LedgerInformation(context.Background())
	if err != nil {
		return 0, fmt.Errorf("[GetBlockNumber] err %v", err)
	}

	number, _ := strconv.Atoi(res.BlockHeight)
	return uint64(number), nil
}

func (c *Apt) Balance(address string) (balance decimal.Decimal, amount string, err error) {
	var res *client.AccountResource
	res, err = c.Client.GetResourceByAccountAddressAndResourceType(
		context.Background(),
		address, "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>",
	)
	if err != nil {
		if strings.Contains(err.Error(), "not_found") {
			d, _ := decimal.NewFromString("0")
			return d, "0", nil
		}
		return
	}
	value := res.Data.Coin.Value
	valueDecimal, _ := decimal.NewFromString(value)
	balance = valueDecimal.Div(decimal.NewFromFloat(math.Pow10(int(c.Decimals))))
	return balance, value, nil
}

func (c *Apt) BalanceOf(address string) (balance decimal.Decimal, err error) {
	if c.ContractAddress == "" {
		return decimal.NewFromFloat(0), nil
	}
	var accountList []client.AccountResource
	accountList, err = c.Client.GetAccountResources(context.Background(), address)
	if err != nil {
		return
	}

	for _, v := range accountList {
		if v.Type == fmt.Sprintf("0x1::coin::CoinStore<%v::coin::T>", c.ContractAddress) {
			value := v.Data.Coin.Value
			_d, _ := decimal.NewFromString(value)
			_d2 := decimal.NewFromFloat(math.Pow10(int(c.Decimals)))
			balance = _d.Div(_d2)
			return
		}
	}

	return decimal.Zero, nil
}

func (c *Apt) GetNonce(address string) (nonce uint64, err error) {
	accountInfo, err := c.Client.GetAccount(context.Background(), address)
	if err != nil {
		return 0, fmt.Errorf("GetNonce err %v", err)
	}
	if accountInfo == nil {
		return 0, nil
	}
	sequenceNumber, _ := strconv.Atoi(accountInfo.SequenceNumber)
	nonce = uint64(sequenceNumber)
	return nonce, nil
}

func (c *Apt) GetGasPrice() (gasPrice *big.Int, err error) {
	res, err := c.Client.EstimateGasPrice(context.Background())
	if err != nil {
		common.Logger.Errorf("get gas price err: %s", err.Error())
		return big.NewInt(0), code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
	}
	gasPrice = big.NewInt(int64(res))
	return gasPrice, nil
}

func (c *Apt) GetGas() (gasLimit decimal.Decimal, gasPrice decimal.Decimal, gasPremium decimal.Decimal, err error) {
	gasrice, err := c.GetGasPrice()
	return c.GasLimit, decimal.NewFromBigInt(gasrice, 0), decimal.Zero, err
}

func (c *Apt) GetTxFee(feeType int64, from, to, amount string, tag string) (fee decimal.Decimal, err error) {
	gasLimit, gasPrice, _, err := c.GetGas()
	if err != nil {
		return
	}
	fee = gasPrice.Mul(gasLimit).Div(decimal.NewFromFloat(math.Pow10(int(c.Decimals))))
	if fee == decimal.Zero {
		fee, _ = decimal.NewFromString(c.MaxFee)
	}
	return
}

func (c *Apt) SendRawTransaction(raw string) (txid string, err error) {
	decodeString, _ := hex.DecodeString(raw)
	tx := new(models.Transaction)
	err = lcs.Unmarshal(decodeString, tx)
	if err != nil {
		return "", err
	}

	transaction, err := c.Client.SubmitTransaction(context.Background(), tx.UserTransaction)

	if err != nil {
		return "", err
	}

	return transaction.Hash, nil
}

func (c *Apt) Sign(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (sig string, err error) {
	fromAddr := CalcAddress(privkey)
	fmt.Printf("from: %s\n", fromAddr)

	nonce, err := c.GetNonce(fromAddr)
	if err != nil {
		common.Logger.Errorf("get nonce err: %s", err.Error())
		if strings.Contains(err.Error(), "account_not_found") {
			return "", code.NewI18nError(code.SrcAccountNotFound, "The sending account does not exist, please check and try again")
		}
		err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
		return
	}

	gasLimit, gasPrice, _, err := c.GetGas()
	if err != nil {
		common.Logger.Errorf("get gas err: %s", err.Error())
		err = code.NewI18nError(code.NetworkErr, "Network error, please try again later.")
		return
	}

	properties := map[string]interface{}{
		"nonce":     strconv.FormatUint(nonce, 10),
		"gas_limit": gasLimit.String(),
		"gas_price": gasPrice.String(),
		"decimals":  "8",
	}

	aptPacker := new(aptPackager)

	err = aptPacker.Pack("dot", properties, coinAddress, amountDec, fromAddr, toAddr)
	if err != nil {
		common.Logger.Errorf("packer err: %s", err.Error())
		return
	}

	tx, err := aptPacker.GetRaw()
	if err != nil {
		err = fmt.Errorf("GetRaw err: %s", err.Error())
		return
	}

	msgBytes, err := tx.GetSigningMessage()
	if err != nil {
		err = fmt.Errorf("GetSigningMessage err: %s", err.Error())
		return
	}

	signature, err := eddsa.Sign(privkey, msgBytes)
	if err != nil {
		return
	}

	err = tx.SetAuthenticator(models.TransactionAuthenticatorEd25519{
		PublicKey: ed25519.PublicKey(eddsa.Pubkey(privkey).Serialize()),
		Signature: signature,
	}).Error()
	if err != nil {
		return
	}

	txBytes, err := lcs.Marshal(tx)
	if err != nil {
		return
	}
	return hex.EncodeToString(txBytes), nil
}

func (c *Apt) Transfer(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (txHash string, err error) {
	sig, err := c.Sign(coinAddress, privkey, toAddr, amountDec)
	if err != nil {
		return
	}
	txHash, err = c.SendRawTransaction(sig)
	if err != nil {
		common.Logger.Errorf("send tx err: %s", err.Error())
		return
	}
	return txHash, nil
}
