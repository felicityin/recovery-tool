package sol

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	solangSdkCommon "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/shopspring/decimal"

	"recovery-tool/tx/eddsa"
	solanaCommon "recovery-tool/tx/sol/solana-adapter/common"
	"recovery-tool/tx/sol/solana-adapter/transaction"
)

const (
	SignTypeAccount = "account"
	SignTypeCreate  = "create"
)

type solanaPackager struct {
	coin     string
	chain    string
	amount   decimal.Decimal
	fromAddr string
	toAddr   string
	tran     *types.Message
	gasLimit string
	gasPrice string
}

func (p *solanaPackager) Pack(
	chain string,
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

	fromTmp, exists := properties["from"]
	if !exists {
		return fmt.Errorf("solana pack fail, from missing")
	}
	from, _ := fromTmp.(string)

	toTmp, exists := properties["to"]
	if !exists {
		return fmt.Errorf("solana pack fail, to missing")
	}
	to, _ := toTmp.(string)

	if from == "" || to == "" {
		return fmt.Errorf("solana pack fail, from or to not empty")
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

	var tokenProgramId string
	tmpTokenProgramId, exists := properties["token_program_id"]
	if !exists || strings.EqualFold(tmpTokenProgramId.(string), "") {
		tokenProgramId = solanaCommon.TokenProgramID.String()
	} else {
		tokenProgramId = tmpTokenProgramId.(string)
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

	var transfer types.Instruction
	var createAccount types.Instruction
	var instructions []types.Instruction
	var toAccount string

	gasPrice := gasPriceInt.Uint64()
	gasLimit := gasLimitInt.Uint64()
	p.gasLimit = tmpGasLimit.(string)
	p.gasPrice = tmpGasPrice.(string)
	if gasPrice > 0 && gasLimit > 0 {
		modifyComputeUnits := transaction.SetComputeUnitLimit(uint32(gasLimit))
		addPriorityFee := transaction.SetComputeUnitPrice(uint64(gasPrice))
		instructions = append(instructions, modifyComputeUnits)
		instructions = append(instructions, addPriorityFee)
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
		createAccount, toAccount, err = transaction.CreateAssociatedTokenAccount(p.fromAddr, p.toAddr, p.coin, tokenProgramId)
		if err != nil {
			return fmt.Errorf("solana pack fail, CreateAssociatedTokenAccount err")
		}
		tp := transaction.TokenTransferParams{
			From:      from,       // sender token address
			To:        toAccount,  // recipient token address
			Authority: p.fromAddr, // signer address
			Decimal:   uint8(decInt),
			Amount:    amountInt,
			Mint:      p.coin, // token mint address
		}
		transfer, err = transaction.NewTokenTransfer(tp, tokenProgramId)
		if err != nil {
			return fmt.Errorf("solana pack fail, NewTokenTransfer err")
		}
		instructions = append(instructions, createAccount)
		instructions = append(instructions, transfer)
	} else if len(p.coin) > 0 {
		tp := transaction.TokenTransferParams{
			From:      from,
			To:        to,
			Authority: p.fromAddr,
			Decimal:   uint8(decInt),
			Amount:    amountInt,
			Mint:      p.coin,
		}
		transfer, err = transaction.NewTokenCheckTransfer(tp, tokenProgramId)
		if err != nil {
			return fmt.Errorf("solana pack fail, NewTokenCheckTransfer err")
		}
		instructions = append(instructions, transfer)
	} else {
		tp := transaction.TransferParams{
			From:   from,
			To:     to,
			Amount: amountInt,
		}
		transfer, err = transaction.NewTransfer(tp)
		if err != nil {
			return fmt.Errorf("solana pack fail, NewTransfer err")
		}
		instructions = append(instructions, transfer)
	}

	blockHash, exists := properties["block_hash"]
	if !exists {
		return fmt.Errorf("solana pack fail, block_hash missing")
	}

	blockHashStr, _ := blockHash.(string)
	if len(blockHashStr) <= 0 {
		return fmt.Errorf("solana pack fail, block_hash not empty")
	}

	feePayer := solangSdkCommon.PublicKeyFromString(p.fromAddr)
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:                   feePayer,
		Instructions:               instructions,
		RecentBlockhash:            blockHashStr,
		AddressLookupTableAccounts: nil,
	})
	p.tran = &message

	return nil
}

func (p *solanaPackager) GetRaw() (*types.Message, error) {
	if p.tran == nil {
		return nil, errors.New("transaction not found")
	}
	return p.tran, nil
}
