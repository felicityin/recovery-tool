package transaction

/*
func：
author： flynn
date: 2020-08-03
fork: https://github.com/solana-labs/solana-web3.js/src/system-program.js
*/
import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	common2 "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/pkg/bincode"
	"github.com/blocto/solana-go-sdk/program/tokenprog"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/btcsuite/btcutil/base58"

	"recovery-tool/tx/sol/solana-adapter/common"
)

type TokenTransferParams struct {
	From      string
	To        string
	Authority string
	Mint      string
	Decimal   uint8
	Amount    *big.Int
}

type TransferParams struct {
	From   string
	To     string
	Amount *big.Int
}

func (tp *TokenTransferParams) GetFromPublicKey() []byte {
	return base58.Decode(tp.From)
}
func (tp *TokenTransferParams) GetToPublicKey() []byte {
	return base58.Decode(tp.To)
}
func (tp *TokenTransferParams) GetAuthPublicKey() []byte {
	return base58.Decode(tp.Authority)
}
func (tp *TokenTransferParams) GetMintPublicKey() []byte {
	return base58.Decode(tp.Mint)
}
func (tp *TransferParams) GetFromPublicKey() []byte {
	return base58.Decode(tp.From)
}

func (tp *TransferParams) GetToPublicKey() []byte {
	return base58.Decode(tp.To)
}
func FindAssociatedTokenAccount(wallet string, tokenMint string, tokenProgramId string) string {

	walletPub := common.PublicKeyFromString(wallet)
	mintPub := common.PublicKeyFromString(tokenMint)
	// signerPub := common.PublicKeyFromString(signer)
	assosiatedAccount, _, err := common.FindAssociatedTokenAddress(walletPub, mintPub, tokenProgramId)
	if err != nil {
		return ""
	}
	return assosiatedAccount.String()
}
func CreateAssociatedTokenAccount(signer string, wallet string, tokenMint string, tokenProgramId string) (types.Instruction, string, error) {

	walletPub := common.PublicKeyFromString(wallet)
	mintPub := common.PublicKeyFromString(tokenMint)
	// signerPub := common.PublicKeyFromString(signer)
	assosiatedAccount, _, err := common.FindAssociatedTokenAddress(walletPub, mintPub, tokenProgramId)
	fmt.Println("assosiatedAccount = ", assosiatedAccount.String())
	if err != nil {
		fmt.Println(err)
	}
	var ti types.Instruction
	ti.Accounts = []types.AccountMeta{
		{PubKey: common2.PublicKeyFromString(signer), IsSigner: true, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(assosiatedAccount.String()), IsSigner: false, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(wallet), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(tokenMint), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(common2.SystemProgramID.String()), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(tokenProgramId), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(common2.SysVarRentPubkey.String()), IsSigner: false, IsWritable: false},
	}
	ti.ProgramID = common2.SPLAssociatedTokenAccountProgramID
	return ti, assosiatedAccount.String(), nil
}

func NewTokenTransfer(transfer TokenTransferParams, tokenProgramId string) (types.Instruction, error) {
	transferIndex := uint8(12) //https://github.com/solana-labs/solana-web3.js/src/system-program.js-->p511  version:v0.64.0
	lamports := transfer.Amount.Uint64()
	decimal := transfer.Decimal

	var ins types.Instruction

	buf1 := new(bytes.Buffer)
	buf2 := make([]byte, 8)
	buf3 := new(bytes.Buffer)

	err := binary.Write(buf1, binary.LittleEndian, transferIndex)
	if err != nil {
		return ins, fmt.Errorf("encode transferIndex error,Err=%v", err)
	}
	binary.LittleEndian.PutUint64(buf2, lamports)

	err = binary.Write(buf3, binary.LittleEndian, decimal)
	if err != nil {
		return ins, fmt.Errorf("encode decimal error,Err=%v", err)
	}
	var data []byte
	data = append(data, buf1.Bytes()...)
	data = append(data, buf2...)

	// if len(data) != 12 {
	// 	return nil, errors.New("transfer data length is not equal 12")
	// }
	data = append(data, buf3.Bytes()...)
	ins.Accounts = []types.AccountMeta{
		{PubKey: common2.PublicKeyFromString(transfer.From), IsSigner: false, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(transfer.Mint), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(transfer.To), IsSigner: false, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(transfer.Authority), IsSigner: true, IsWritable: true},
	}

	ins.ProgramID = common2.PublicKeyFromString(tokenProgramId)
	ins.Data = data
	return ins, nil
}

func NewTokenCheckTransfer(transfer TokenTransferParams, tokenProgramId string) (types.Instruction, error) {
	transferIndex := uint8(12) //https://github.com/solana-labs/solana-web3.js/src/system-program.js-->p511  version:v0.64.0
	lamports := transfer.Amount.Uint64()
	decimal := transfer.Decimal

	buf1 := new(bytes.Buffer)
	buf2 := make([]byte, 8)
	buf3 := new(bytes.Buffer)

	var ins types.Instruction

	err := binary.Write(buf1, binary.LittleEndian, transferIndex)
	if err != nil {
		return ins, fmt.Errorf("encode transferIndex error,Err=%v", err)
	}
	binary.LittleEndian.PutUint64(buf2, lamports)

	err = binary.Write(buf3, binary.LittleEndian, decimal)
	if err != nil {
		return ins, fmt.Errorf("encode decimal error,Err=%v", err)
	}
	var data []byte
	data = append(data, buf1.Bytes()...)
	data = append(data, buf2...)

	// if len(data) != 12 {
	// 	return nil, errors.New("transfer data length is not equal 12")
	// }
	data = append(data, buf3.Bytes()...)

	ins.Accounts = []types.AccountMeta{
		{PubKey: common2.PublicKeyFromString(transfer.From), IsSigner: false, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(transfer.Mint), IsSigner: false, IsWritable: false},
		{PubKey: common2.PublicKeyFromString(transfer.To), IsSigner: false, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(transfer.Authority), IsSigner: true, IsWritable: true},
	}
	ins.ProgramID = common2.PublicKeyFromString(tokenProgramId)
	ins.Data = data
	return ins, nil
}

func NewTokenCheckTransferV2(transfer TokenTransferParams, tokenProgramId string) types.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction tokenprog.Instruction
		Amount      uint64
		Decimals    uint8
	}{
		Instruction: tokenprog.InstructionTransferChecked,
		Amount:      transfer.Amount.Uint64(),
		Decimals:    transfer.Decimal,
	})
	if err != nil {
		panic(err)
	}

	accounts := make([]types.AccountMeta, 0, 4)
	accounts = append(accounts, types.AccountMeta{PubKey: common2.PublicKeyFromString(transfer.From), IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: common2.PublicKeyFromString(transfer.Mint), IsSigner: false, IsWritable: false})
	accounts = append(accounts, types.AccountMeta{PubKey: common2.PublicKeyFromString(transfer.To), IsSigner: false, IsWritable: true})
	accounts = append(accounts, types.AccountMeta{PubKey: common2.PublicKeyFromString(transfer.Authority), IsSigner: true, IsWritable: false})
	//for _, signerPubkey := range param.Signers {
	//	accounts = append(accounts, types.AccountMeta{PubKey: signerPubkey, IsSigner: true, IsWritable: false})
	//}

	return types.Instruction{
		ProgramID: common2.PublicKeyFromString(tokenProgramId),
		Accounts:  accounts,
		Data:      data,
	}
}

func NewTransfer(transfer TransferParams) (types.Instruction, error) {
	transferIndex := uint32(2) //https://github.com/solana-labs/solana-web3.js/src/system-program.js-->p511  version:v0.64.0
	lamports := transfer.Amount.Uint64()
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	var ins types.Instruction
	err := binary.Write(buf1, binary.LittleEndian, transferIndex)
	if err != nil {
		return ins, fmt.Errorf("encode transfer index error,Err=%v", err)
	}
	err = binary.Write(buf2, binary.LittleEndian, lamports)
	if err != nil {
		return ins, fmt.Errorf("encode lamports error,Err=%v", err)
	}
	var data []byte
	data = append(data, buf1.Bytes()...)
	data = append(data, buf2.Bytes()...)
	if len(data) != 12 {
		return ins, errors.New("transfer data length is not equal 12")
	}

	ins.Accounts = []types.AccountMeta{
		{PubKey: common2.PublicKeyFromString(transfer.From), IsSigner: true, IsWritable: true},
		{PubKey: common2.PublicKeyFromString(transfer.To), IsSigner: false, IsWritable: true},
	}
	ins.ProgramID = common2.SystemProgramID
	ins.Data = data
	return ins, nil
}
