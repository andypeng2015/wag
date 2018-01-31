// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x86

import (
	"github.com/tsavola/wag/internal/gen"
	"github.com/tsavola/wag/internal/links"
	"github.com/tsavola/wag/internal/opers"
	"github.com/tsavola/wag/internal/regs"
	"github.com/tsavola/wag/internal/values"
	"github.com/tsavola/wag/traps"
	"github.com/tsavola/wag/types"
)

func (mach X86) BinaryOp(code gen.RegCoder, oper uint16, a, b values.Operand) values.Operand {
	if (oper & opers.BinaryFloat) == 0 {
		switch {
		case (oper & opers.BinaryCompare) != 0:
			return mach.binaryIntCompareOp(code, uint8(oper), a, b)

		case (oper & opers.BinaryIntShift) != 0:
			return mach.binaryIntShiftOp(code, uint8(oper), a, b)

		case (oper & opers.BinaryIntDivmul) != 0:
			return mach.binaryIntDivmulOp(code, uint8(oper), a, b)

		default:
			return mach.commonBinaryIntOp(code, uint8(oper), a, b)
		}
	} else {
		switch {
		case (oper & opers.BinaryCompare) != 0:
			return mach.binaryFloatCompareOp(code, uint8(oper), a, b)

		case (oper & opers.BinaryFloatMinmax) != 0:
			return mach.binaryFloatMinmaxOp(code, uint8(oper), a, b)

		case (oper & opers.BinaryFloatCopysign) != 0:
			return mach.binaryFloatCopysignOp(code, a, b)

		default:
			return mach.commonBinaryFloatOp(code, uint8(oper), a, b)
		}
	}
}

var commonBinaryIntInsns = []binaryInsn{
	opers.IndexIntAdd: Add,
	opers.IndexIntSub: Sub,
	opers.IndexIntAnd: And,
	opers.IndexIntOr:  Or,
	opers.IndexIntXor: Xor,
}

func (mach X86) commonBinaryIntOp(code gen.RegCoder, index uint8, a, b values.Operand) (result values.Operand) {
	if index == opers.IndexIntSub && a.Storage == values.Imm && a.ImmValue() == 0 {
		return mach.inplaceIntOp(code, Neg, b)
	}

	switch b.Storage {
	case values.Imm:
		switch {
		case b.ImmValue() == 1:
			switch index {
			case opers.IndexIntAdd:
				return mach.inplaceIntOp(code, Inc, a)

			case opers.IndexIntSub:
				return mach.inplaceIntOp(code, Dec, a)
			}

		case b.ImmValue() < -0x80000000 || b.ImmValue() >= 0x80000000:
			b = mach.opBorrowMaybeScratchRegOperand(code, b, true)
		}

	case values.VarReference, values.Stack, values.ConditionFlags:
		b = mach.opBorrowMaybeScratchRegOperand(code, b, true)
	}

	insn := commonBinaryIntInsns[index]

	targetReg, _ := mach.opMaybeResultReg(code, a, false)
	result = values.TempRegOperand(a.Type, targetReg, true)

	switch {
	case b.Storage.IsReg():
		insn.opFromReg(code, a.Type, targetReg, b.Reg())
		code.Consumed(b)
		return

	case b.Storage == values.VarMem:
		insn.opFromStack(code, a.Type, targetReg, b.VarMemOffset())
		return

	case b.Storage == values.Imm: // large values moved to registers earlier
		insn.opImm(code, a.Type, targetReg, int32(b.ImmValue()))
		return

	default:
		panic("unexpected storage type of second operand of common binary int op")
	}
}

func (mach X86) binaryIntCompareOp(code gen.RegCoder, cond uint8, a, b values.Operand) (result values.Operand) {
	targetReg, _, own := mach.opBorrowMaybeResultReg(code, a, false)
	if own {
		defer code.FreeReg(a.Type, targetReg)
	}

	result = values.ConditionFlagsOperand(values.Condition(cond))

	switch {
	case b.Storage.IsReg():
		Cmp.opFromReg(code, a.Type, targetReg, b.Reg())
		code.Consumed(b)
		return

	case b.Storage == values.VarMem:
		Cmp.opFromStack(code, a.Type, targetReg, b.VarMemOffset())
		return

	case b.Storage == values.Imm && b.ImmValue() >= -0x80000000 && b.ImmValue() < 0x80000000:
		Cmp.opImm(code, a.Type, targetReg, int32(b.ImmValue()))
		return

	default:
		mach.OpMove(code, regScratch, b, false)
		Cmp.opFromReg(code, a.Type, targetReg, regScratch)
		return
	}
}

var binaryDivmulInsns = []struct {
	insnRexM
	shiftImm shiftImmInsn
}{
	opers.IndexDivmulDivS: {Idiv, NoShiftImmInsn},
	opers.IndexDivmulDivU: {Div, ShrImm},
	opers.IndexDivmulRemS: {Idiv, NoShiftImmInsn},
	opers.IndexDivmulRemU: {Div, NoShiftImmInsn}, // TODO: use AND for 2^n divisors
	opers.IndexDivmulMul:  {Mul, ShlImm},
}

func (mach X86) binaryIntDivmulOp(code gen.RegCoder, index uint8, a, b values.Operand) values.Operand {
	insn := binaryDivmulInsns[index]
	t := a.Type

	if b.Storage == values.Imm {
		value := b.ImmValue()

		switch {
		case insn.shiftImm.defined() && value > 0 && isPowerOfTwo(uint64(value)):
			reg, _ := mach.opMaybeResultReg(code, a, false)
			insn.shiftImm.op(code, t, reg, log2(uint64(value)))
			return values.TempRegOperand(t, reg, true)
		}
	}

	division := (index & opers.DivmulMul) == 0
	checkZero := true
	checkOverflow := true

	if b.Storage.IsReg() {
		if b.Reg() == regResult {
			newReg := regScratch

			if division {
				var ok bool

				// can't use scratch reg as divisor since it contains the dividend high bits
				newReg, ok = code.TryAllocReg(t)
				if !ok {
					// borrow a register which we don't need in this function
					MovMMX.opFromReg(code, types.I64, regScratchMMX, regTextBase)
					defer MovMMX.opToReg(code, types.I64, regTextBase, regScratchMMX)

					newReg = regTextBase
				}
			}

			Mov.opFromReg(code, t, newReg, regResult)
			b = values.RegOperand(true, t, newReg)
		}
	} else {
		if division && b.Storage == values.Imm {
			value := b.ImmValue()
			if value != 0 {
				checkZero = false
			}
			if value != -1 {
				checkOverflow = false
			}
		}

		reg, ok := code.TryAllocReg(t)
		if !ok {
			// borrow a register which we don't need in this function
			MovMMX.opFromReg(code, types.I64, regScratchMMX, regTextBase)
			defer MovMMX.opToReg(code, types.I64, regTextBase, regScratchMMX)

			reg = regTextBase
		}

		mach.OpMove(code, reg, b, true)
		b = values.RegOperand(true, t, reg)
	}

	mach.OpMove(code, regResult, a, false)

	remainder := (index & opers.DivmulRem) != 0

	var doNot links.L

	if division {
		if checkZero {
			mach.opCheckDivideByZero(code, t, b.Reg())
		}

		if a.Storage == values.Imm {
			value := a.ImmValue()
			if t.Size() == types.Size32 {
				if value != -0x80000000 {
					checkOverflow = false
				}
			} else {
				if value != -0x8000000000000000 {
					checkOverflow = false
				}
			}
		}

		signed := (index & opers.DivmulSign) != 0

		if signed && checkOverflow {
			var do links.L

			if remainder {
				Xor.opFromReg(code, types.I32, regScratch, regScratch) // moved to result at the end

				Cmp.opImm(code, t, b.Reg(), -1)
				Je.rel8.opStub(code)
				doNot.AddSite(code.Len())
			} else {
				switch t.Size() {
				case types.Size32:
					Cmp.opImm(code, t, regResult, -0x80000000)

				case types.Size64:
					MovImm64.op(code, t, regScratch, -0x8000000000000000)
					Cmp.opFromReg(code, t, regResult, regScratch)

				default:
					panic(a)
				}

				Jne.rel8.opStub(code)
				do.AddSite(code.Len())

				Cmp.opImm(code, t, b.Reg(), -1)
				Jne.rel8.opStub(code)
				do.AddSite(code.Len())

				code.OpTrapCall(traps.IntegerOverflow)
			}

			do.Addr = code.Len()
			mach.updateBranches8(code, &do)
		}

		if signed {
			// sign-extend dividend low bits to high bits
			CdqCqo.op(code, t)
		} else {
			// zero-extend dividend high bits
			Xor.opFromReg(code, types.I32, regScratch, regScratch)
		}
	}

	insn.opReg(code, t, b.Reg())
	code.Consumed(b)

	doNot.Addr = code.Len()
	mach.updateBranches8(code, &doNot)

	if remainder {
		Mov.opFromReg(code, t, regResult, regScratch)
	}

	return values.TempRegOperand(t, regResult, true)
}

func (mach X86) opCheckDivideByZero(code gen.RegCoder, t types.T, reg regs.R) {
	var end links.L

	Test.opFromReg(code, t, reg, reg)
	Jne.rel8.opStub(code)
	end.AddSite(code.Len())

	code.OpTrapCall(traps.IntegerDivideByZero)

	end.Addr = code.Len()
	mach.updateBranches8(code, &end)
}

var binaryShiftInsns = []struct {
	reg insnRexM
	imm shiftImmInsn
}{
	opers.IndexShiftRotl: {Rol, RolImm},
	opers.IndexShiftRotr: {Ror, RorImm},
	opers.IndexShiftShl:  {Shl, ShlImm},
	opers.IndexShiftShrS: {Sar, SarImm},
	opers.IndexShiftShrU: {Shr, ShrImm},
}

func (mach X86) binaryIntShiftOp(code gen.RegCoder, index uint8, a, b values.Operand) (result values.Operand) {
	insn := binaryShiftInsns[index]

	switch {
	case b.Storage == values.Imm:
		reg, _ := mach.opMaybeResultReg(code, a, true)
		insn.imm.op(code, b.Type, reg, uint8(b.ImmValue()))
		result = values.TempRegOperand(a.Type, reg, true)

	case b.Storage.IsReg() && b.Reg() == regShiftCount:
		reg, _ := mach.opMaybeResultReg(code, a, false)
		insn.reg.opReg(code, a.Type, reg)
		code.Discard(b)
		result = values.TempRegOperand(a.Type, reg, true)

	case code.RegAllocated(types.I32, regShiftCount):
		reg, _ := mach.opMaybeResultReg(code, a, true)
		if reg == regShiftCount {
			Mov.opFromReg(code, a.Type, regResult, regShiftCount)
			result = mach.subtleShiftOp(code, insn.reg, a.Type, regResult, b)
			code.FreeReg(types.I32, regShiftCount)
		} else {
			// unknown operand in regShiftCount
			Mov.opFromReg(code, types.I64, regScratch, regShiftCount) // save
			result = mach.subtleShiftOp(code, insn.reg, a.Type, reg, b)
			Mov.opFromReg(code, types.I64, regShiftCount, regScratch) // restore
		}

	default:
		code.AllocSpecificReg(types.I32, regShiftCount)
		reg, _ := mach.opMaybeResultReg(code, a, true)
		result = mach.subtleShiftOp(code, insn.reg, a.Type, reg, b)
		code.FreeReg(types.I32, regShiftCount)
	}

	return
}

// subtleShiftOp trashes regShiftCount.
func (mach X86) subtleShiftOp(code gen.RegCoder, insn insnRexM, t types.T, reg regs.R, count values.Operand) values.Operand {
	count.Type = types.I32                         // TODO: 8-bit mov
	mach.OpMove(code, regShiftCount, count, false) //
	insn.opReg(code, t, reg)
	return values.TempRegOperand(t, reg, true)
}

var commonBinaryFloatInsns = []insnPrefix{
	opers.IndexFloatAdd: AddsSSE,
	opers.IndexFloatSub: SubsSSE,
	opers.IndexFloatDiv: DivsSSE,
	opers.IndexFloatMul: MulsSSE,
}

// TODO: support memory source operands

func (mach X86) commonBinaryFloatOp(code gen.RegCoder, index uint8, a, b values.Operand) values.Operand {
	targetReg, _ := mach.opMaybeResultReg(code, a, false)

	sourceReg, _, own := mach.opBorrowMaybeScratchReg(code, b, false)
	if own {
		defer code.FreeReg(b.Type, sourceReg)
	}

	commonBinaryFloatInsns[index].opFromReg(code, a.Type, targetReg, sourceReg)
	return values.TempRegOperand(a.Type, targetReg, false)
}

var binaryFloatMinmaxInsns = []struct {
	commonInsn insnPrefix
	zeroInsn   insnPrefix
}{
	opers.IndexMinmaxMin: {MinsSSE, OrpSSE},
	opers.IndexMinmaxMax: {MaxsSSE, AndpSSE},
}

func (mach X86) binaryFloatMinmaxOp(code gen.RegCoder, index uint8, a, b values.Operand) values.Operand {
	targetReg, _ := mach.opMaybeResultReg(code, a, false)

	sourceReg, _, own := mach.opBorrowMaybeScratchReg(code, b, false)
	if own {
		defer code.FreeReg(b.Type, sourceReg)
	}

	var common links.L
	var end links.L

	UcomisSSE.opFromReg(code, a.Type, targetReg, sourceReg)
	Jne.rel8.opStub(code)
	common.AddSite(code.Len())

	binaryFloatMinmaxInsns[index].zeroInsn.opFromReg(code, a.Type, targetReg, sourceReg)
	JmpRel.rel8.opStub(code)
	end.AddSite(code.Len())

	common.Addr = code.Len()
	mach.updateBranches8(code, &common)

	binaryFloatMinmaxInsns[index].commonInsn.opFromReg(code, a.Type, targetReg, sourceReg)

	end.Addr = code.Len()
	mach.updateBranches8(code, &end)

	return values.TempRegOperand(a.Type, targetReg, false)
}

func (mach X86) binaryFloatCompareOp(code gen.RegCoder, cond uint8, a, b values.Operand) values.Operand {
	aReg, _, own := mach.opBorrowMaybeResultReg(code, a, true)
	if own {
		defer code.FreeReg(a.Type, aReg)
	}

	bReg, _, own := mach.opBorrowMaybeScratchReg(code, b, false)
	if own {
		defer code.FreeReg(b.Type, bReg)
	}

	UcomisSSE.opFromReg(code, a.Type, aReg, bReg)
	return values.ConditionFlagsOperand(values.Condition(cond))
}

func (mach X86) binaryFloatCopysignOp(code gen.RegCoder, a, b values.Operand) values.Operand {
	targetReg, _ := mach.opMaybeResultReg(code, a, false)

	sourceReg, _, own := mach.opBorrowMaybeScratchReg(code, b, false)
	if own {
		defer code.FreeReg(b.Type, sourceReg)
	}

	panic("TODO")

	return values.TempRegOperand(a.Type, targetReg, false)
}
