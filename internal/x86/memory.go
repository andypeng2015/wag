// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/internal/gen"
	"github.com/tsavola/wag/internal/links"
	"github.com/tsavola/wag/internal/opers"
	"github.com/tsavola/wag/internal/regs"
	"github.com/tsavola/wag/internal/values"
	"github.com/tsavola/wag/trap"
	"github.com/tsavola/wag/wasm"
)

type memoryAccess struct {
	insn     binaryInsn
	insnType abi.Type
	zeroExt  bool
}

var memoryLoads = []memoryAccess{
	opers.IndexIntLoad:    {mov, 0, true},
	opers.IndexIntLoad8S:  {binaryInsn{movsx8, noPrefixMIInsn}, 0, false},
	opers.IndexIntLoad8U:  {binaryInsn{movzx8, noPrefixMIInsn}, 0, false},
	opers.IndexIntLoad16S: {binaryInsn{movsx16, noPrefixMIInsn}, 0, false},
	opers.IndexIntLoad16U: {binaryInsn{movzx16, noPrefixMIInsn}, 0, false},
	opers.IndexIntLoad32S: {binaryInsn{movsxd, noPrefixMIInsn}, 0, false}, // type is ignored
	opers.IndexIntLoad32U: {mov, abi.I32, true},
	opers.IndexFloatLoad:  {binaryInsn{movsSSE, noPrefixMIInsn}, 0, false},
}

var memoryStores = []memoryAccess{
	opers.IndexIntStore:   {mov, 0, false},
	opers.IndexIntStore8:  {mov8, abi.I32, false},
	opers.IndexIntStore16: {mov16, abi.I32, false},
	opers.IndexIntStore32: {mov, abi.I32, false},
	opers.IndexFloatStore: {binaryInsn{movsSSE, movImm}, 0, false}, // integer immediate works
}

// LoadOp makes sure that index gets zero-extended if it's a VarReg operand.
func (ISA) LoadOp(m *Module, code gen.RegCoder, oper uint16, index values.Operand, resultType abi.Type, offset uint32) (result values.Operand) {
	size := oper >> 8

	baseReg, indexReg, ownIndexReg, disp := opMemoryAddress(m, code, size, index, offset)
	if ownIndexReg {
		defer code.FreeReg(abi.I64, indexReg)
	}

	load := memoryLoads[uint8(oper)]

	targetReg, ok := code.TryAllocReg(resultType)
	if !ok {
		targetReg = RegResult
	}

	result = values.TempRegOperand(resultType, targetReg, load.zeroExt)

	insnType := load.insnType
	if insnType == 0 {
		insnType = resultType
	}

	load.insn.opFromIndirect(&m.Text, insnType, targetReg, 0, indexReg, baseReg, disp)
	return
}

// StoreOp makes sure that index gets zero-extended if it's a VarReg operand.
func (ISA) StoreOp(m *Module, code gen.RegCoder, oper uint16, index, x values.Operand, offset uint32) {
	size := oper >> 8

	baseReg, indexReg, ownIndexReg, disp := opMemoryAddress(m, code, size, index, offset)
	if ownIndexReg {
		defer code.FreeReg(abi.I64, indexReg)
	}

	store := memoryStores[uint8(oper)]

	insnType := store.insnType
	if insnType == 0 {
		insnType = x.Type
	}

	if x.Storage == values.Imm {
		value := x.ImmValue()
		value32 := int32(value)

		switch {
		case size == 1:
			value32 = int32(int8(value32))

		case size == 2:
			value32 = int32(int16(value32))

		case size == 4 || (value >= -0x80000000 && value < 0x80000000):

		default:
			goto large
		}

		store.insn.opImmToIndirect(&m.Text, insnType, 0, indexReg, baseReg, disp, value32)
		return

	large:
	}

	valueReg, _, own := opBorrowMaybeResultReg(m, code, x, false)
	if own {
		defer code.FreeReg(x.Type, valueReg)
	}

	store.insn.opToIndirect(&m.Text, insnType, valueReg, 0, indexReg, baseReg, disp)
}

// opMemoryAddress may return the scratch register as the base.
func opMemoryAddress(m *Module, code gen.Coder, size uint16, index values.Operand, offset uint32) (baseReg, indexReg regs.R, ownIndexReg bool, disp int32) {
	sizeReach := uint64(size - 1)
	reachOffset := uint64(offset) + sizeReach

	if reachOffset >= 0x80000000 {
		code.OpTrapCall(trap.MemoryOutOfBounds)
		return
	}

	alreadyChecked := reachOffset < uint64(index.Bounds.Upper)

	switch index.Storage {
	case values.Imm:
		value := uint64(index.ImmValue())

		if value >= 0x80000000 {
			code.OpTrapCall(trap.MemoryOutOfBounds)
			return
		}

		addr := value + uint64(offset)
		reachAddr := addr + sizeReach

		if reachAddr >= 0x80000000 {
			code.OpTrapCall(trap.MemoryOutOfBounds)
			return
		}

		if reachAddr < uint64(code.MinMemorySize()) || alreadyChecked {
			baseReg = RegMemoryBase
			indexReg = NoIndex
			disp = int32(addr)
			return
		}

		lea.opFromIndirect(&m.Text, abi.I64, RegScratch, 0, NoIndex, RegMemoryBase, int32(reachAddr))

	default:
		reg, zeroExt, own := opBorrowMaybeScratchReg(m, code, index, true)

		if !zeroExt {
			mov.opFromReg(&m.Text, abi.I32, reg, reg) // zero-extend index
		}

		if alreadyChecked {
			baseReg = RegMemoryBase
			indexReg = reg
			ownIndexReg = own
			disp = int32(offset)
			return
		}

		lea.opFromIndirect(&m.Text, abi.I64, RegScratch, 0, reg, RegMemoryBase, int32(reachOffset))

		if own {
			code.FreeReg(abi.I32, reg)
		}
	}

	cmp.opFromReg(&m.Text, abi.I64, RegScratch, RegMemoryLimit)

	if addr := code.TrapTrampolineAddr(trap.MemoryOutOfBounds); addr != 0 {
		jge.op(&m.Text, addr)
	} else {
		var checked links.L

		jl.rel8.opStub(&m.Text)
		checked.AddSite(m.Text.Pos())

		code.OpTrapCall(trap.MemoryOutOfBounds)

		checked.Addr = m.Text.Pos()
		updateLocalBranches(m, &checked)
	}

	baseReg = RegScratch
	indexReg = NoIndex
	disp = -int32(sizeReach)
	return
}

func (ISA) OpCurrentMemory(m *Module) values.Operand {
	mov.opFromReg(&m.Text, abi.I64, RegResult, RegMemoryLimit)
	sub.opFromReg(&m.Text, abi.I64, RegResult, RegMemoryBase)
	shrImm.op(&m.Text, abi.I64, RegResult, wasm.PageBits)

	return values.TempRegOperand(abi.I32, RegResult, true)
}

func (ISA) OpGrowMemory(m *Module, code gen.RegCoder, x values.Operand) values.Operand {
	var out links.L
	var fail links.L

	movMMX.opToReg(&m.Text, abi.I64, RegScratch, RegMemoryGrowLimitMMX)

	targetReg, zeroExt := opMaybeResultReg(m, code, x, false)
	if !zeroExt {
		mov.opFromReg(&m.Text, abi.I32, targetReg, targetReg)
	}

	shlImm.op(&m.Text, abi.I64, targetReg, wasm.PageBits)
	add.opFromReg(&m.Text, abi.I64, targetReg, RegMemoryLimit) // new memory limit
	cmp.opFromReg(&m.Text, abi.I64, targetReg, RegScratch)

	jg.rel8.opStub(&m.Text)
	fail.AddSite(m.Text.Pos())

	mov.opFromReg(&m.Text, abi.I64, RegScratch, RegMemoryLimit)
	mov.opFromReg(&m.Text, abi.I64, RegMemoryLimit, targetReg)
	sub.opFromReg(&m.Text, abi.I64, RegScratch, RegMemoryBase)
	shrImm.op(&m.Text, abi.I64, RegScratch, wasm.PageBits) // value on success
	mov.opFromReg(&m.Text, abi.I32, targetReg, RegScratch)

	jmpRel.rel8.opStub(&m.Text)
	out.AddSite(m.Text.Pos())

	fail.Addr = m.Text.Pos()
	updateLocalBranches(m, &fail)

	movImm.opImm(&m.Text, abi.I32, targetReg, -1) // value on failure

	out.Addr = m.Text.Pos()
	updateLocalBranches(m, &out)

	return values.TempRegOperand(abi.I32, targetReg, true)
}
