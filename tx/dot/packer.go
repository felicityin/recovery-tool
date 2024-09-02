package dot

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"

	"recovery-tool/common/code"
	"recovery-tool/tx/dot/polkadot-adapter/address"
	"recovery-tool/tx/eddsa"

	"recovery-tool/tx/dot/polkadot-adapter/polkadotTransaction"
)

type dotPackager struct {
	coin     string
	chain    string
	amount   decimal.Decimal
	fromAddr string
	toAddr   string
	tran     *polkadotTransaction.TxStruct
}

func (p *dotPackager) Pack(
	chain string,
	properties map[string]interface{},
	coin string,
	amount decimal.Decimal,
	fromAddr string,
	toAddr string,
) error {
	p.chain = chain
	p.coin = coin
	p.amount = amount
	p.fromAddr = fromAddr
	p.toAddr = toAddr

	senderPubKey, err := address.AddressDecode(p.fromAddr)
	if err != nil {
		return code.NewI18nError(code.InvalidPrivkey, "Private key is invalid")
	}
	recipientPubKey, err := address.AddressDecode(p.toAddr)
	if err != nil {
		return code.NewI18nError(code.InvalidToAddr, "Dest address is invalid")
	}

	nonce, exists := properties["nonce"]
	if !exists {
		return fmt.Errorf("dot pack fail, nonce missing")
	}
	nonceInt := new(big.Int)
	var ok bool
	nonceInt, ok = nonceInt.SetString(nonce.(string), 10)
	if !ok {
		return fmt.Errorf("dot pack fail, convert nonce to big.int err")
	}

	fee, exists := properties["fee"]
	if !exists {
		return fmt.Errorf("dot pack fail, fee missing")
	}

	feeFloat, err := strconv.ParseFloat(fee.(string), 64)
	if err != nil {
		return fmt.Errorf("dot pack fail, convert fee str to float fail")
	}

	blockHash, exists := properties["block_hash"]
	if !exists {
		return fmt.Errorf("dot pack fail, block_hash missing")
	}

	genesisHash, exists := properties["genesis_hash"]
	if !exists {
		return fmt.Errorf("dot pack fail, genesis_hash missing")
	}

	specVersion, exists := properties["spec_version"]
	if !exists {
		return fmt.Errorf("dot pack fail, spec_version missing")
	}

	txVersion, exists := properties["transaction_version"]
	if !exists {
		return fmt.Errorf("dot pack fail, transaction_version missing")
	}
	specVersionInt, _ := strconv.Atoi(specVersion.(string))
	txVersionInt, _ := strconv.Atoi(txVersion.(string))

	dec, exists := properties["decimals"]
	if !exists {
		return fmt.Errorf("dot pack fail, decimals missing")
	}
	decimalsInt, err := strconv.Atoi(dec.(string))
	if err != nil {
		return fmt.Errorf("dot pack fail, convert decimals to int err")
	}

	feeVal := new(big.Float).Mul(big.NewFloat(feeFloat), big.NewFloat(math.Pow10(decimalsInt)))
	feeUint64, _ := feeVal.Uint64()

	amountInt, err := eddsa.BigMulDecimal(p.amount, decimalsInt)
	if err != nil {
		return fmt.Errorf("dot pack fail, amount mul err")
	}

	GenesisHash := RemoveOxToAddress(genesisHash.(string))
	p.tran = &polkadotTransaction.TxStruct{
		//发送方公钥
		SenderPubkey: hex.EncodeToString(senderPubKey),
		//接收方公钥
		RecipientPubkey: hex.EncodeToString(recipientPubKey),
		//发送金额（最小单位）
		Amount: amountInt.Uint64(),
		//nonce
		Nonce: nonceInt.Uint64(),
		//手续费（最小单位）
		Fee: feeUint64,
		//当前高度
		//BlockHeight: uint64(blockHeightInt),
		//当前高度区块哈希
		BlockHash: RemoveOxToAddress(blockHash.(string)),
		//创世块哈希
		GenesisHash: RemoveOxToAddress(GenesisHash),

		// POLKADOT
		// GenesisHash: RemoveOxToAddress("0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3"),
		// KUSAMA
		// GenesisHash: RemoveOxToAddress("b0a8d493285c2df73290dfb7e61f870f17b41801197a149ca93654499ea3dafe"),
		//spec版本
		SpecVersion: uint32(specVersionInt),
		//Transaction版本
		TxVersion: uint32(txVersionInt),
	}
	return nil
}

func (p *dotPackager) GetRaw() (*polkadotTransaction.TxStruct, error) {
	if p.tran == nil {
		return nil, errors.New("transaction not found")
	}
	return p.tran, nil
}

func GetTransferCode() string {
	transferCode := polkadotTransaction.DOT_Balannce_Transfer
	return transferCode
}

func RemoveOxToAddress(addr string) string {
	if strings.Index(addr, "0x") == 0 {
		return addr[2:]
	}
	return addr
}
