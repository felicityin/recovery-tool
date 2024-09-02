package apt

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/portto/aptos-go-sdk/models"
	"github.com/shopspring/decimal"

	"recovery-tool/common/code"
	"recovery-tool/tx/eddsa"
)

type aptPackager struct {
	coin     string
	chain    string
	amount   decimal.Decimal
	fromAddr string
	toAddr   string
	tran     *models.Transaction
}

func (p *aptPackager) Pack(
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

	chainId := 1 // MainNetChain
	nonce, exists := properties["nonce"]
	if !exists {
		return fmt.Errorf("apt pack fail, nonce missing")
	}

	gasLimit, exists := properties["gas_limit"]
	if !exists {
		return fmt.Errorf("apt pack fail, gas_limit missing")
	}
	gasLimitInt := new(big.Int)
	gasLimitInt, ok := gasLimitInt.SetString(gasLimit.(string), 10)
	if !ok {
		return fmt.Errorf("apt pack fail, convert gas_limit to big.int err")
	}

	gasPrice, exists := properties["gas_price"]
	if !exists {
		return fmt.Errorf("apt pack fail, gas_price missing")
	}
	gasPriceInt := new(big.Int)
	gasPriceInt, ok = gasPriceInt.SetString(gasPrice.(string), 10)
	if !ok {
		return fmt.Errorf("apt pack fail, convert gas_price to big.int err")
	}

	dec, exists := properties["decimals"]
	if !exists {
		return fmt.Errorf("apt pack fail, decimals missing")
	}
	decimalsInt, err := strconv.Atoi(dec.(string))
	if err != nil {
		return fmt.Errorf("apt pack fail, convert decimals to int err")
	}

	addr0x1, err := models.HexToAccountAddress("0x1")
	if err != nil {
		return fmt.Errorf("apt pack fail, get addr0x1 err")
	}

	toAcc, err := models.HexToAccountAddress(p.toAddr)
	if err != nil {
		return code.NewI18nError(code.InvalidToAddr, "Dest address is invalid")
	}

	amountInt, err := eddsa.BigMulDecimal(p.amount, decimalsInt)
	if err != nil {
		return fmt.Errorf("apt pack fail, amount mul err")
	}

	p.tran = &models.Transaction{}
	err = p.tran.SetChainID(uint8(chainId)).
		SetSender(p.fromAddr).
		SetPayload(models.EntryFunctionPayload{
			Module: models.Module{
				Address: addr0x1,
				Name:    "aptos_account",
			},
			Function: "transfer",
			Arguments: []interface{}{
				toAcc,
				amountInt.Uint64(),
			},
			TypeArguments: []models.TypeTag{},
		},
		).SetExpirationTimestampSecs(uint64(time.Now().Add(10 * time.Minute).Unix())).
		SetGasUnitPrice(gasPriceInt.Uint64()).
		SetMaxGasAmount(gasLimitInt.Uint64()).
		SetSequenceNumber(nonce).Error()
	if err != nil {
		return err
	}
	return nil
}

func (p *aptPackager) GetRaw() (*models.Transaction, error) {
	if p.tran == nil {
		return nil, errors.New("GetRaw transaction not found")
	}
	return p.tran, nil
}
