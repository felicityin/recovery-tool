package ton

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"github.com/xssnick/tonutils-go/address"
	"github.com/xssnick/tonutils-go/tlb"
	"github.com/xssnick/tonutils-go/ton/jetton"
	"github.com/xssnick/tonutils-go/ton/wallet"
	"github.com/xssnick/tonutils-go/tvm/cell"

	"recovery-tool/common"
	"recovery-tool/tx/eddsa"
)

type tonPackager struct {
	coin      string
	chain     string
	amount    decimal.Decimal
	fromAddr  string
	toAddr    string
	tran      *cell.Cell
	notActive bool
}

// status -1 重播  0 新交易 1 加速
// properties包含字段
// decimals
// token_program_id
// type = account create
// gas_limit gas_price
// block_hash
func (p *tonPackager) Pack(
	chain string,
	properties map[string]interface{},
	coin string,
	amount decimal.Decimal,
	fromAddr, toAddr string,
) error {
	p.chain = chain
	p.coin = coin
	p.amount = amount
	p.fromAddr = fromAddr // 原始地址
	p.toAddr = toAddr     // 原始地址

	dec, exists := properties["decimals"]
	if !exists {
		return fmt.Errorf("ton pack fail, decimals missing")
	}
	decInt := dec.(int)

	amountInt, err := eddsa.BigMulDecimal(p.amount, decInt)
	if err != nil {
		return fmt.Errorf("ton pack fail, amount mul err")
	}

	feeStr, exists := properties["fee"]
	if !exists {
		return fmt.Errorf("ton pack fail, fee missing")
	}
	feeFloat, err := strconv.ParseFloat(feeStr.(string), 64)
	if err != nil {
		return fmt.Errorf("ton pack fail, convert fee str to float fail")
	}
	feeVal := new(big.Float).Mul(big.NewFloat(feeFloat), big.NewFloat(math.Pow10(9)))
	fee64, _ := feeVal.Int64()
	common.Logger.Infof("fee: %d", fee64)

	nonce, exists := properties["nonce"]
	if !exists {
		return fmt.Errorf("ton pack fail, nonce missing")
	}
	nonceInt := nonce.(uint64)

	tmpIsActive, exists := properties["is_active"]
	if !exists {
		return fmt.Errorf("solana pack fail, isActive missing")
	}
	isActive := tmpIsActive.(bool)

	tmpStatus, exists := properties["status"]
	if !exists {
		return fmt.Errorf("solana pack fail, status missing")
	}
	status1 := tmpStatus.(string)

	p.notActive = !isActive || status1 != tlb.AccountStatusActive

	tmpMemoStr, exists := properties["memo"]
	if !exists {
		return fmt.Errorf("solana pack fail, memo missing")
	}
	memoStr := tmpMemoStr.(string)

	var memo *cell.Cell
	if memoStr != "" {
		memo, err = wallet.CreateCommentCell(memoStr)
		if err != nil {
			return fmt.Errorf("ton pack faile, CreateCommentCell err: %s", err.Error())
		}
	}
	messages := make([]*wallet.Message, 0)

	if p.coin == "" {
		message := &wallet.Message{
			Mode: 1, // pay fees separately (from balance, not from amount)
			InternalMessage: &tlb.InternalMessage{
				IHRDisabled: true,
				Bounce:      false, // return amount in case of processing error
				Bounced:     true,
				DstAddr:     address.MustParseAddr(p.toAddr),
				Amount:      tlb.FromNanoTON(amountInt),
				Body:        memo,
			},
		}
		messages = append(messages, message)
	} else {
		tmpJettonAddres, exists := properties["jetton_address"]
		if !exists {
			return fmt.Errorf("solana pack fail, jetton_addres missing")
		}
		jettonAddress := tmpJettonAddres.(string)

		tmpForwardfee, exists := properties["forward_fee"]
		if !exists {
			return fmt.Errorf("solana pack fail, forward_fee missing")
		}
		forwardFee := tmpForwardfee.(int64)

		buf := make([]byte, 8)
		if _, err := rand.Read(buf); err != nil {
			return err
		}
		rnd := binary.LittleEndian.Uint64(buf)
		to := address.MustParseAddr(p.toAddr)
		var commentCell *cell.Cell
		commentCell, err = wallet.CreateCommentCell(memoStr)
		if err != nil {
			return fmt.Errorf("ton pack faile, CreateCommentCell err: %s", err.Error())
		}
		body, err1 := tlb.ToCell(jetton.TransferPayload{
			QueryID:     rnd,
			Amount:      tlb.FromNanoTON(amountInt),
			Destination: to,
			//address where to send a response with confirmation of a successful burn and the rest of the incoming message coins.
			ResponseDestination: address.MustParseAddr(p.fromAddr),
			CustomPayload:       nil,
			ForwardTONAmount:    tlb.FromNanoTON(new(big.Int).SetInt64(forwardFee)),
			ForwardPayload:      commentCell,
		})
		if err1 != nil {
			return fmt.Errorf("ton pack faile, ToCell err: %s", err1.Error())
		}
		message := &wallet.Message{
			Mode: 1,
			InternalMessage: &tlb.InternalMessage{
				IHRDisabled: true,
				Bounce:      false,
				//jetton wallet address
				DstAddr: address.MustParseAddr(jettonAddress),
				Amount:  tlb.FromNanoTON(new(big.Int).SetInt64(fee64)),
				Body:    body,
			},
		}
		messages = append(messages, message)
	}

	if len(messages) > 4 {
		return fmt.Errorf("ton pack fail, for this type of wallet max 4 messages can be sent in the same time")
	}

	payload := cell.BeginCell().MustStoreUInt(uint64(wallet.DefaultSubwallet), 32).
		MustStoreUInt(uint64(time.Now().Add(time.Duration(1000)*time.Second).UTC().Unix()), 32).
		MustStoreUInt(nonceInt, 32)

	for i, m := range messages {
		intMsg, err := tlb.ToCell(m.InternalMessage)
		if err != nil {
			return fmt.Errorf("ton pack fail, failed to convert internal message %d to cell: %w", i, err)
		}
		payload.MustStoreUInt(uint64(m.Mode), 8).MustStoreRef(intMsg)
	}

	p.tran = payload.EndCell()
	return nil
}

func (p *tonPackager) GetRaw() (*cell.Cell, error) {
	if p.tran == nil {
		return nil, errors.New("transaction not found")
	}
	return p.tran, nil
}

func (p *tonPackager) NotActive() bool {
	return p.notActive
}
