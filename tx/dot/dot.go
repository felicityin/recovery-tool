package dot

import (
	"encoding/hex"
	"fmt"
	"math"
	"strconv"

	"recovery-tool/tx/dot/polkadot-adapter/polkadotTransaction"

	"github.com/mr-tron/base58"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/blake2b"

	"recovery-tool/common"
	"recovery-tool/tx/eddsa"
)

var ss58Prefix = []byte("SS58PRE")
var DOTNetWorkByteMap = map[string]byte{
	"DOT": 0x00,
	"KSM": 0x00,
}

func CalcAddress(privkey []byte) (string, error) {
	publicKey := eddsa.Pubkey(privkey)
	return DOTPublicKeyToAddress(publicKey.Serialize(), DOTNetWorkByteMap["DOT"])
}

func DOTPublicKeyToAddress(pub []byte, network byte) (string, error) {
	enc := append([]byte{network}, pub...)
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return "", err
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return "", err
	}
	checksum := hasher.Sum(nil)

	return base58.Encode(append(enc, checksum[:2]...)), nil
}

type Dot struct {
	Client      *Client
	GenesisHash string
	Decimals    int
}

func NewDot(url string) *Dot {
	client := NewClient(url)

	dot := new(Dot)
	dot.Client = client
	dot.GenesisHash = "0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3"
	dot.Decimals = 10
	return dot
}

type BlockInfo struct {
	BlockHash          string
	BlockHeight        string
	SpecVersion        string
	TransactionVersion string
	GenesisHash        string
}

func (c *Dot) GetBlockInfo() (blockInfo BlockInfo, err error) {
	var runTimeRes RuntimeSpecResponse
	runTimeRes, err = c.Client.RuntimeSpec()
	if err != nil {
		return BlockInfo{}, err
	}

	var blocksHeadResponse BlocksHeadResponse
	blocksHeadResponse, err = c.Client.BlocksHead()
	if err != nil {
		return BlockInfo{}, err
	}

	blockInfo.BlockHash = blocksHeadResponse.Hash
	blockInfo.BlockHeight = blocksHeadResponse.Number
	blockInfo.SpecVersion = runTimeRes.SpecVersion
	blockInfo.TransactionVersion = runTimeRes.TransactionVersion
	blockInfo.GenesisHash = c.GenesisHash
	return
}

func (c *Dot) GetBlockNumber() (uint64, error) {
	res, err := c.Client.BlocksHead()
	if err != nil {
		return 0, fmt.Errorf("[GetBlockNumber] err %v", err)
	}

	number, _ := strconv.Atoi(res.Number)
	return uint64(number), nil
}

func (c *Dot) Balance(address string) (balance decimal.Decimal, amount string, err error) {
	var res AccountsBalanceInfoResponse
	res, err = c.Client.AccountsBalanceInfo(address)

	if err != nil {
		return
	}

	value := res.Free
	valueDecimal, _ := decimal.NewFromString(value)
	balance = valueDecimal.Div(decimal.NewFromFloat(math.Pow10(int(c.Decimals))))
	return balance, value, nil
}

func (c *Dot) SendRawTransaction(raw string) (txid string, err error) {
	var transaction TransactionResponse
	transaction, err = c.Client.Transaction(raw)
	txid = transaction.Hash
	return
}

func (c *Dot) GetNonce(address string) (nonce uint64, err error) {
	var res AccountsBalanceInfoResponse
	res, err = c.Client.AccountsBalanceInfo(address)
	if err != nil {
		return
	}
	nonceInt, _ := strconv.Atoi(res.Nonce)
	nonce = uint64(nonceInt)
	return nonce, nil
}

func (c *Dot) Sign(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (sig string, err error) {
	blockInfo, err := c.GetBlockInfo()
	if err != nil {
		return
	}

	fromAddr, err := CalcAddress(privkey)
	if err != nil {
		return "", err
	}
	common.Logger.Infof("from: %s\n", fromAddr)

	nonce, err := c.GetNonce(fromAddr)
	if err != nil {
		return "", fmt.Errorf("get nonce err: %s", err.Error())
	}

	properties := map[string]interface{}{
		"nonce":               strconv.FormatUint(nonce, 10),
		"fee":                 "0",
		"decimals":            "10",
		"block_hash":          blockInfo.BlockHash,
		"block_height":        blockInfo.BlockHeight,
		"genesis_hash":        blockInfo.GenesisHash,
		"spec_version":        blockInfo.SpecVersion,
		"transaction_version": blockInfo.TransactionVersion,
	}

	dotPacker := new(dotPackager)

	err = dotPacker.Pack("dot", properties, coinAddress, amountDec, fromAddr, toAddr)
	if err != nil {
		err = fmt.Errorf("packer err: %s", err.Error())
		return
	}

	txRaw, err := dotPacker.GetRaw()
	if err != nil {
		err = fmt.Errorf("GetRaw err: %s", err.Error())
		return
	}

	msg, err := txRaw.NewTxPayLoad(GetTransferCode())
	if err != nil {
		return "", err
	}

	msgBytes, err := hex.DecodeString(msg.NewVersionToBytesString())
	if err != nil {
		err = fmt.Errorf("GetSignPacket err: %s", err.Error())
		return
	}

	signature, err := eddsa.Sign(privkey, msgBytes)
	if err != nil {
		return
	}

	tx, pass := polkadotTransaction.VerifyAndCombineTransaction(
		GetTransferCode(),
		txRaw.ToJSONString(),
		hex.EncodeToString(signature),
	)
	if !pass {
		return "", fmt.Errorf("polkadot verify failed")
	}
	return tx, nil
}

func (c *Dot) Transfer(coinAddress string, privkey []byte, toAddr string, amountDec decimal.Decimal) (txHash string, err error) {
	sig, err := c.Sign(coinAddress, privkey, toAddr, amountDec)
	if err != nil {
		return
	}
	txHash, err = c.SendRawTransaction(sig)
	if err != nil {
		return
	}
	return txHash, nil
}
