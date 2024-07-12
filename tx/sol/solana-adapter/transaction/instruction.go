/**
 * Created by Goland.
 * User: Jett
 * Date: 2024/3/17
 * Time: 8:36 AM
 */
package transaction

import (
	"encoding/binary"
	common2 "github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/near/borsh-go"
)

type Instruction borsh.Enum

const (
	InstructionRequestUnits uint8 = iota
	InstructionRequestHeapFrame
	InstructionSetComputeUnitLimit
	InstructionSetComputeUnitPrice
)

type RequestUnitsParam struct {
	Units         uint32
	AdditionalFee uint32
}

// RequestUnits ...
//func RequestUnits(param RequestUnitsParam) TransactionInstruction {
//	data, err := borsh.Serialize(struct {
//		Instruction   Instruction
//		Units         uint32
//		AdditionalFee uint32
//	}{
//		Instruction:   InstructionRequestUnits,
//		Units:         param.Units,
//		AdditionalFee: param.AdditionalFee,
//	})
//	if err != nil {
//		panic(err)
//	}
//	return TransactionInstruction{
//		programId: common.ComputeBudgetProgramID.String(),
//		keys:      []*AccountMeta{},
//		data:      data,
//	}
//}

type RequestHeapFrameParam struct {
	Bytes uint32
}

// RequestHeapFrame ...
//func RequestHeapFrame(param RequestHeapFrameParam) types.Instruction {
//	data, err := borsh.Serialize(struct {
//		Instruction Instruction
//		Bytes       uint32
//	}{
//		Instruction: InstructionRequestHeapFrame,
//		Bytes:       param.Bytes,
//	})
//	if err != nil {
//		panic(err)
//	}
//
//	return types.Instruction{
//		ProgramID: common.ComputeBudgetProgramID,
//		Accounts:  []types.AccountMeta{},
//		Data:      data,
//	}
//}

type SetComputeUnitLimitParam struct {
	Units uint32
}

// SetComputeUnitLimit set a specific compute unit limit that the transaction is allowed to consume.
func SetComputeUnitLimit(units uint32) types.Instruction {
	data := make([]byte, 1+4)
	data[0] = InstructionSetComputeUnitLimit
	binary.LittleEndian.PutUint32(data[1:], units)
	//return &TransactionInstruction{
	//	programId: common.ComputeBudgetProgramID.String(),
	//	//keys:      []*AccountMeta{},
	//	data: data,
	//}

	var ti types.Instruction
	ti.ProgramID = common2.ComputeBudgetProgramID
	ti.Data = data
	//ti.SetKeys([]*AccountMeta{})
	return ti
}

type SetComputeUnitPriceParam struct {
	MicroLamports uint64
}

// SetComputeUnitPrice set a compute unit price in "micro-lamports" to pay a higher transaction
// fee for higher transaction prioritization.
func SetComputeUnitPrice(microLamports uint64) types.Instruction {
	data := make([]byte, 1+8)
	data[0] = InstructionSetComputeUnitPrice
	binary.LittleEndian.PutUint64(data[1:], microLamports)

	//return &TransactionInstruction{
	//	programId: common.ComputeBudgetProgramID.String(),
	//	//keys:      []*AccountMeta{},
	//	data: data,
	//}

	var ti types.Instruction
	ti.ProgramID = common2.ComputeBudgetProgramID
	ti.Data = data
	//ti.SetKeys([]*AccountMeta{})
	return ti
}
