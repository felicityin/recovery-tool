package sol

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	solangSdkCommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/assotokenprog"
	"github.com/blocto/solana-go-sdk/program/cmptbdgprog"
	"github.com/blocto/solana-go-sdk/program/sysprog"
	"github.com/blocto/solana-go-sdk/program/tokenprog"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/shopspring/decimal"

	"recovery-tool/tx/eddsa"
)

type solanaPackagerV2 struct {
	coin     string
	chain    string
	amount   decimal.Decimal
	fromAddr string
	toAddr   string
	tran     *types.Message
	gasLimit string
	gasPrice string
}

func (p *solanaPackagerV2) Pack(chain string,
	properties map[string]interface{},
	coin string,
	amount decimal.Decimal,
	fromAddr, toAddr string,
) error {
	p.chain = chain
	p.coin = coin
	p.amount = amount
	p.fromAddr = fromAddr
	p.toAddr = toAddr

	feePayer := solangSdkCommon.PublicKeyFromString(p.fromAddr)
	destPub := solangSdkCommon.PublicKeyFromString(p.toAddr)

	fromTmp, exists := properties["from"]
	if !exists {
		return fmt.Errorf("solana pack fail, from address missing")
	}
	from, _ := fromTmp.(string)

	toTmp, exists := properties["to"]
	if !exists {
		return fmt.Errorf("solana pack fail, to address missing")
	}
	to, _ := toTmp.(string)

	if from == "" || to == "" {
		return fmt.Errorf("solana pack fail, from or to address not empty")
	}

	var decInt int
	var err error

	dec, exists := properties["decimals"]
	if !exists {
		return fmt.Errorf("solana pack fail, decimals missing")
	}

	decInt, err = strconv.Atoi(dec.(string))
	if err != nil {
		return fmt.Errorf("solana pack fail, decimals to int err")
	}

	if p.coin != "" {
		if decInt <= 0 {
			return fmt.Errorf("solana pack fail, decimals should be > 0")
		}
	}

	tmpGasLimit, exists := properties["gas_limit"]
	if !exists {
		return fmt.Errorf("solana pack fail, gas_limit missing")
	}

	gasLimitInt := new(big.Int)
	gasLimitInt, ok := gasLimitInt.SetString(tmpGasLimit.(string), 10)
	if !ok {
		return fmt.Errorf("solana pack fail, convert gas_limit to big.int err")
	}

	tmpGasPrice, exists := properties["gas_price"]
	if !exists {
		return fmt.Errorf("solana pack fail, gas_price missing")
	}
	gasPriceInt := new(big.Int)
	gasPriceInt, ok = gasPriceInt.SetString(tmpGasPrice.(string), 10)
	if !ok {
		return fmt.Errorf("solana pack fail, convert gas_price to big.int err")
	}

	instructions := make([]types.Instruction, 0)

	gasPrice := gasPriceInt.Uint64()
	gasLimit := gasLimitInt.Uint64()
	p.gasLimit = tmpGasLimit.(string)
	p.gasPrice = tmpGasPrice.(string)
	if gasPrice > 0 && gasLimit > 0 {
		instructions = append(instructions, cmptbdgprog.SetComputeUnitPrice(cmptbdgprog.SetComputeUnitPriceParam{MicroLamports: gasPrice}))
		instructions = append(instructions, cmptbdgprog.SetComputeUnitLimit(cmptbdgprog.SetComputeUnitLimitParam{Units: uint32(gasLimit)}))
	}

	amountInt, err := eddsa.BigMulDecimal(p.amount, decInt)
	if err != nil {
		return fmt.Errorf("solana pack fail, amount mul err")
	}

	var signType string
	tmpType, exists := properties["type"]
	if exists {
		signType = tmpType.(string)
	}

	if signType == SignTypeAccount {
		contractToken := solangSdkCommon.PublicKeyFromString(p.coin)
		toAccount, _, err := solangSdkCommon.FindAssociatedTokenAddress(destPub, contractToken)
		if err != nil {
			return fmt.Errorf("solana pack fail, FindAssociatedTokenAddress err: %s", err.Error())
		}

		instructions = append(instructions, assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
			Funder:                 feePayer,
			Owner:                  destPub,
			Mint:                   contractToken,
			AssociatedTokenAccount: toAccount,
		}))

		instructions = append(instructions, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
			From:     solangSdkCommon.PublicKeyFromString(from),
			To:       toAccount,
			Auth:     feePayer,
			Signers:  []solangSdkCommon.PublicKey{},
			Amount:   amountInt.Uint64(),
			Decimals: uint8(decInt),
			Mint:     contractToken,
		}))
	} else if len(p.coin) > 0 {
		if from == p.fromAddr || to == p.toAddr {
			return fmt.Errorf("solana pack fail, token address invalid")
		}
		contractToken := solangSdkCommon.PublicKeyFromString(p.coin)

		instructions = append(instructions, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
			From:     solangSdkCommon.PublicKeyFromString(from),
			To:       solangSdkCommon.PublicKeyFromString(to),
			Auth:     feePayer,
			Signers:  []solangSdkCommon.PublicKey{feePayer},
			Amount:   amountInt.Uint64(),
			Decimals: uint8(decInt),
			Mint:     contractToken,
		}))
	} else {
		instructions = append(instructions, sysprog.Transfer(sysprog.TransferParam{
			From:   feePayer,
			To:     destPub,
			Amount: amountInt.Uint64(),
		}))
	}

	blockHash, exists := properties["block_hash"]
	if !exists {
		return fmt.Errorf("solana pack fail, block_hash missing")
	}
	blockHashStr, _ := blockHash.(string)

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   feePayer,
		Instructions:               instructions,
		RecentBlockhash:            blockHashStr,
		AddressLookupTableAccounts: nil,
	})
	p.tran = &message

	return nil
}

func (p *solanaPackagerV2) GetRaw() (*types.Message, error) {
	if p.tran == nil {
		return nil, errors.New("transaction not found")
	}
	return p.tran, nil
}

func (p *solanaPackagerV2) GetSignPacket() (string, error) {
	msgByte, err := p.tran.Serialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(msgByte), nil
}
