package ton

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"

	"github.com/ipfs/go-log"
	"github.com/shopspring/decimal"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/liteclient"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/nft"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"

	cm "recovery-tool/common"
	"recovery-tool/tx/eddsa"
)

const Usdt = "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
const TONForwardFee = 10000000
const MaxFee = "0.050000"

type Ton struct {
	Client           *Client
	ApiClient        ton.APIClientWrapped
	tokenAccount     string
	Decimals         int
	ContractAddress  string
	ContractDecimals int
	log              *log.ZapEventLogger
}

func NewTon(url string) *Ton {
	client := liteclient.NewConnectionPool()

	// connect to testnet lite server
	err := client.AddConnectionsFromConfigUrl(context.Background(), "https://ton.org/global.config.json")
	if err != nil {
		panic(err)
	}

	// initialize ton api lite connection wrapper
	api := ton.NewAPIClient(client).WithRetry()

	ton := &Ton{
		Client:    NewClient(url),
		ApiClient: api,
		Decimals:  9,
		log:       cm.Logger,
	}
	return ton
}

func (this *Ton) Balance(address string) (string, decimal.Decimal, error) {
	res, err := this.Client.Balance(address)
	if err != nil {
		this.log.Errorf("balance failed: %v", err)
		return "", decimal.Zero, err
	}
	amount, _ := decimal.NewFromString(res.Balance)
	amount = amount.Div(decimal.NewFromFloat(math.Pow10(int(this.Decimals))))
	return res.Balance, amount, nil
}

func (this *Ton) BalanceOf(address string) (decimals int, amount string, balance decimal.Decimal, err error) {
	this.ContractDecimals, err = this.GetTokenDeciamals()
	if err != nil {
		this.log.Errorf("GetTokenDeciamals err: %s", err.Error())
		return
	}

	if this.tokenAccount == "" {
		this.tokenAccount, err = this.GetTokenAddr(address)
		if err != nil {
			this.log.Errorf("GetTonTokenAddress err: %s", err.Error())
			return
		}
	}
	this.log.Infof("token address: %s", this.tokenAccount)

	jettonBal, err := this.Client.JettonBalance(address, this.tokenAccount, this.ContractAddress)
	if err != nil {
		this.log.Errorf("get getton balance err: %s", err.Error())
		return
	}
	if len(jettonBal.JettonWallets) == 0 {
		return 0, amount, balance, fmt.Errorf("no jettonBalance result, %s-%s-%s", address, this.ContractAddress, this.tokenAccount)
	}
	balance, err = decimal.NewFromString(jettonBal.JettonWallets[0].Balance)
	balance = balance.Div(decimal.NewFromFloat(math.Pow10(int(this.ContractDecimals))))
	return this.ContractDecimals, jettonBal.JettonWallets[0].Balance, balance, nil
}

func (c *Ton) Transfer(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (txHash string, err error) {
	pubkey := eddsa.Pubkey(privkey).Serialize()

	from, err := wallet.AddressFromPubKey(pubkey, wallet.V3, wallet.DefaultSubwallet)
	if err != nil {
		c.log.Errorf("gen addr err: %s", err.Error())
		return
	}
	fromAddr := from.String()
	c.log.Infof("from: %s", fromAddr)

	account, err := c.GetAccount(fromAddr)
	if err != nil {
		c.log.Errorf("GetAccount err: %s", err.Error())
		return
	}

	nonce, err := c.GetNonce(fromAddr)
	if err != nil {
		c.log.Errorf("GetNonce err: %s", err.Error())
		return
	}
	cm.Logger.Infof("nonce: %d", nonce)

	fee := c.GetTxFee()

	properties := map[string]interface{}{
		"decimals":  c.Decimals,
		"fee":       fee.String(),
		"nonce":     nonce,
		"is_active": account.IsActive,
		"status":    account.Status,
		"memo":      "",
	}

	if coinAddress != "" {
		c.ContractAddress = coinAddress
		c.ContractDecimals, err = c.GetTokenDeciamals()
		if err != nil {
			c.log.Errorf("GetTokenDeciamals err: %", err.Error())
			return
		}

		c.tokenAccount, err = c.GetTokenAddr(fromAddr)
		if err != nil {
			c.log.Errorf("GetTonTokenAddress err: %s", err.Error())
			return
		}

		properties["decimals"] = c.ContractDecimals
		properties["jetton_address"] = c.tokenAccount
		properties["forward_fee"] = int64(TONForwardFee)
	}

	tonPacker := new(tonPackager)

	err = tonPacker.Pack("ton", properties, coinAddress, amountDec, fromAddr, toAddr)
	if err != nil {
		c.log.Errorf("pack err: %s", err.Error())
		return
	}

	tran, err := tonPacker.GetRaw()
	if err != nil {
		c.log.Errorf("packer.GetRaw err: %s", err.Error())
		return
	}

	signature, err := eddsa.Sign(privkey, tran.Hash())
	if err != nil {
		c.log.Errorf("sign err: %s", err.Error())
		return
	}

	if !tran.Verify(pubkey, signature) {
		c.log.Errorf("verify sig failed")
		return "", fmt.Errorf("verify sig failed")
	}

	msg := cell.BeginCell().MustStoreSlice(signature, 512).MustStoreBuilder(tran.ToBuilder()).EndCell()

	var stateInit *tlb.StateInit
	if tonPacker.NotActive() {
		stateInit, err = wallet.GetStateInit(pubkey, wallet.V3, wallet.DefaultSubwallet)
		if err != nil {
			c.log.Errorf("failed to get init state: %s", err.Error())
			return
		}
	}

	externalMessage := &tlb.ExternalMessage{
		DstAddr:   address.MustParseAddr(fromAddr),
		StateInit: stateInit,
		Body:      msg,
	}

	cell, err := tlb.ToCell(externalMessage)
	if err != nil {
		c.log.Errorf("tlb.ToCell err: %s", err.Error())
		return
	}

	marshalJSON, err := cell.MarshalJSON()
	if err != nil {
		c.log.Errorf("cell.MarshalJSON err: %s", err.Error())
		return
	}
	rawTx := base64.StdEncoding.EncodeToString(marshalJSON)

	txHash, err = c.SendRawTransaction(rawTx)
	if err != nil {
		c.log.Errorf("send tx err: %s", err.Error())
		return
	}
	return txHash, nil
}

func (this *Ton) GetTokenDeciamals() (int, error) {
	tokenContract := address.MustParseAddr(this.ContractAddress)
	master := jetton.NewJettonMasterClient(this.ApiClient, tokenContract)

	data, err := master.GetJettonData(context.Background())
	if err != nil {
		this.log.Errorf("GetJettonData err: %s", err.Error())
		return 0, err
	}

	content := data.Content.(*nft.ContentSemichain).ContentOnchain
	decimals, err := strconv.Atoi(content.GetAttribute("decimals"))
	if err != nil {
		this.log.Errorf("invalid decimals: %s", err.Error())
	}
	return decimals, nil
}

func (this *Ton) GetTokenAddr(addr string) (string, error) {
	tokenContract := address.MustParseAddr(this.ContractAddress)
	master := jetton.NewJettonMasterClient(this.ApiClient, tokenContract)

	tokenWallet, err := master.GetJettonWallet(context.Background(), address.MustParseAddr(addr))
	if err != nil {
		this.log.Errorf("GetJettonWallet err: %s", err.Error())
		return "", err
	}

	return tokenWallet.Address().String(), nil
}

type AccountInfo struct {
	IsActive bool
	Status   string
}

func (this *Ton) GetAccount(name string) (account AccountInfo, err error) {
	res, err := this.Client.Balance(name)
	if err != nil {
		this.log.Errorf("balance failed: %v", err)
		return account, err
	}
	return AccountInfo{
		IsActive: res.IsActive,
		Status:   res.Status,
	}, nil
}

func (this *Ton) GetNonce(address string) (nonce uint64, err error) {
	nonce, err = this.Client.Sequence(address)
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (this Ton) SendRawTransaction(raw string) (txid string, err error) {
	res, err := this.Client.BroadCast(raw)
	if err != nil {
		this.log.Errorf("broadcast tx failed: %s", err.Error())
		if res != nil {
			this.log.Errorf("sendRawTransaction failed: %+v", res)
		}
		return "", err
	}
	if res.MessageHash == "" {
		this.log.Errorf("sendRawTransaction failed: %+v", res)
		return "", fmt.Errorf("sendRawTransaction failed: %+v", res)
	}
	return res.MessageHash, err
}

func (this *Ton) GetTxFee() decimal.Decimal {
	fee, _ := decimal.NewFromString(MaxFee)
	amt := decimal.NewFromInt(TONForwardFee).Div(decimal.NewFromFloat(math.Pow10(int(this.Decimals))))
	return fee.Add(amt)
}
