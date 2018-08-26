// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package codegen

import (
	"errors"
	"fmt"

	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/internal/gen"
	"github.com/tsavola/wag/internal/gen/link"
	"github.com/tsavola/wag/internal/gen/reg"
	"github.com/tsavola/wag/internal/gen/val"
	"github.com/tsavola/wag/internal/loader"
	"github.com/tsavola/wag/internal/obj"
	"github.com/tsavola/wag/internal/typedecode"
)

func pushBranchTarget(f *gen.Func, valueType abi.Type, funcEnd bool) {
	stackOffset := f.StackOffset

	if int(f.NumInitedVars) < len(f.Vars) {
		// init still in progress, but any branch expressions will have
		// initialized all vars before we reach the target
		numRegParams := int32(len(f.Vars)) - f.NumStackParams
		stackOffset = numRegParams * obj.Word
	}

	f.BranchTargets = append(f.BranchTargets, &gen.BranchTarget{
		StackOffset: stackOffset,
		ValueType:   valueType,
		FuncEnd:     funcEnd,
	})
}

func popBranchTarget(f *gen.Func) (finalizedLabel *link.L) {
	n := len(f.BranchTargets) - 1
	finalizedLabel = &f.BranchTargets[n].Label
	f.BranchTargets = f.BranchTargets[:n]

	trimBoundsStacks(f)
	return
}

func getBranchTarget(f *gen.Func, depth uint32) *gen.BranchTarget {
	if depth >= uint32(len(f.BranchTargets)) {
		panic(fmt.Errorf("relative branch depth out of bounds: %d", depth))
	}
	return f.BranchTargets[len(f.BranchTargets)-int(depth)-1]
}

func boundsStackLevel(f *gen.Func) int {
	return len(f.BranchTargets)
}

func trimBoundsStacks(f *gen.Func) {
	size := boundsStackLevel(f) + 1
	for i := range f.Vars {
		f.Vars[i].TrimBoundsStack(size)
	}
}

func genBlock(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	t := typedecode.Block(load.Varint7())

	opSaveTemporaryOperands(f)
	opInitVars(f)
	opStoreVars(f, false)

	pushBranchTarget(f, t, false) // end

	savedMinBlockOperand := f.MinBlockOperand
	f.MinBlockOperand = len(f.Operands)

	deadend = genOps(f, load)

	var result val.Operand

	if deadend {
		for len(f.Operands) > f.MinBlockOperand {
			x := popOperand(f)
			debugf("discarding operand at end of %s due to deadend: %s", op, x)
			discard(f, x)
		}
	} else {
		if t != abi.Void {
			result = popOperand(f)
			if result.Type != t {
				panic(fmt.Errorf("%s result has wrong type: %s", op, result.Type))
			}
		}

		if len(f.Operands) != f.MinBlockOperand {
			panic(fmt.Errorf("operands remain on stack after %s", op))
		}
	}

	f.MinBlockOperand = savedMinBlockOperand

	if end := popBranchTarget(f); end.Live() {
		if result.Storage != val.Nowhere {
			opMove(f, reg.Result, result, false)
		}

		if t != abi.Void {
			result = val.TempRegOperand(t, reg.Result, false)
		}

		opLabel(f, end)
		isa.UpdateBranches(f.Text.Bytes(), end)
		deadend = false
	}

	if result.Storage != val.Nowhere {
		pushOperand(f, result)
	}

	return
}

func genBr(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	relativeDepth := load.Varuint32()
	target := getBranchTarget(f, relativeDepth)

	if target.ValueType != abi.Void {
		value := popOperand(f)
		if value.Type != target.ValueType {
			panic(fmt.Errorf("%s value operand type is %s, but target expects %s", op, value.Type, target.ValueType))
		}
		opMove(f, reg.Result, value, false)
	}

	if target.FuncEnd {
		asm.AddStackPtrImm(f.M, f.StackOffset)
		asm.Return(f.M)
	} else {
		opSaveTemporaryOperands(f) // TODO: avoid saving operands which we are going to skip over
		opInitVars(f)
		opStoreVars(f, true)
		asm.AddStackPtrImm(f.M, f.StackOffset-target.StackOffset)
		opBranch(f, &target.Label)
	}

	deadend = true
	return
}

func genBrIf(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	relativeDepth := load.Varuint32()
	target := getBranchTarget(f, relativeDepth)

	cond := opPreloadOperand(f, popOperand(f))
	if cond.Type != abi.I32 {
		panic(fmt.Errorf("%s: condition operand has wrong type: %s", op, cond.Type))
	}

	var value val.Operand

	if target.ValueType != abi.Void {
		value = popOperand(f)
		if value.Type != target.ValueType {
			panic(fmt.Errorf("%s: value operand has wrong type: %s", op, value.Type))
		}
	}

	opSaveTemporaryOperands(f)
	opInitVars(f)
	opStoreVars(f, false)

	if value.Type != abi.Void {
		if cond.Storage == val.TempReg && cond.Reg() == reg.Result {
			r := opAllocReg(f, abi.I32)
			asm.MoveI32Reg(f.M, r, reg.Result)
			cond = val.TempRegOperand(cond.Type, r, true)
		}

		opMove(f, reg.Result, value, true)
	}

	stackDelta := f.StackOffset - target.StackOffset

	asm.AddStackPtrImm(f.M, stackDelta)
	opBranchIf(f, cond, true, &target.Label)
	asm.AddStackPtrImm(f.M, -stackDelta)

	if target.ValueType != abi.Void {
		pushResultRegOperand(f, target.ValueType)
	}
	return
}

func genBrTable(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	targetCount := load.Varuint32()
	if targetCount >= uint32(MaxBranchTableSize) {
		panic(fmt.Errorf("%s has too many targets: %d", op, targetCount))
	}

	targetTable := make([]*gen.BranchTarget, targetCount)

	for i := range targetTable {
		relativeDepth := load.Varuint32()
		target := getBranchTarget(f, relativeDepth)
		target.Label.SetLive()
		targetTable[i] = target
	}

	relativeDepth := load.Varuint32()
	defaultTarget := getBranchTarget(f, relativeDepth)
	defaultTarget.Label.SetLive()

	index := opPreloadOperand(f, popOperand(f))
	if index.Type != abi.I32 {
		panic(fmt.Errorf("%s: index operand has wrong type: %s", op, index.Type))
	}

	valueType := defaultTarget.ValueType

	for i, target := range targetTable {
		if target.ValueType != valueType {
			panic(fmt.Errorf("%s targets have inconsistent value types: %s (default target) vs. %s (target #%d)", op, valueType, target.ValueType, i))
		}
	}

	var value val.Operand

	if valueType != abi.Void {
		value = popOperand(f)
		if value.Type != valueType {
			panic(fmt.Errorf("%s: value operand has wrong type: %s", op, value.Type))
		}
	}

	var commonStackOffset int32
	var tableType = abi.I32
	var tableScale uint8 = 2

	if len(targetTable) > 0 {
		commonStackOffset = targetTable[0].StackOffset
		for _, target := range targetTable[1:] {
			if target.StackOffset != commonStackOffset {
				commonStackOffset = -1
				tableType = abi.I64
				tableScale = 3
				break
			}
		}
	}

	tableSize := len(targetTable) << tableScale
	alignMask := (1 << tableScale) - 1
	tableAddr := (len(f.ROData.Bytes()) + alignMask) &^ alignMask
	f.ROData.ResizeBytes(tableAddr + tableSize)

	opSaveTemporaryOperands(f) // TODO: avoid saving operands which we are going to skip over?
	opInitVars(f)
	opStoreVars(f, false)

	if value.Type != abi.Void {
		if index.Storage == val.TempReg && index.Reg() == reg.Result {
			r := opAllocReg(f, abi.I32)
			asm.MoveI32Reg(f.M, r, reg.Result)
			index = val.TempRegOperand(index.Type, r, true)
		}

		opMove(f, reg.Result, value, true)
	}

	var r reg.R

	if index.Storage.IsReg() {
		r = index.Reg()
	} else {
		r = opAllocReg(f, abi.I32)
		asm.Move(f, r, index, false)
	}

	f.Regs.FreeAll()

	// vars were already stored and reg freed
	for i := range f.Vars {
		f.Vars[i].ResetCache()
	}

	// if index yielded a var register, then it was just freed, but the
	// register retains its value.  don't call anything that allocates
	// registers until the critical section ends.

	defaultDelta := f.StackOffset - defaultTarget.StackOffset

	asm.AddStackPtrImm(f.M, defaultDelta)
	tableStackOffset := f.StackOffset - defaultDelta
	opBranchIfOutOfBounds(f, r, int32(len(targetTable)), &defaultTarget.Label)

	asm.LoadIntROData(f, tableType, r, tableScale, int32(tableAddr))

	var zeroExt bool

	if commonStackOffset >= 0 {
		asm.AddStackPtrImm(f.M, tableStackOffset-commonStackOffset)
		zeroExt = true
	} else {
		asm.AddStackPtrUpper32(f.M, r)
		zeroExt = false
	}

	asm.BranchIndirect(f.M, r, zeroExt)

	// end of critical section.

	t := gen.BranchTable{
		RODataAddr: int32(tableAddr),
		Targets:    targetTable,
	}
	if commonStackOffset >= 0 {
		t.CodeStackOffset = -1
	} else {
		// no common offset
		t.CodeStackOffset = tableStackOffset
	}
	f.BranchTables = append(f.BranchTables, t)

	deadend = true
	return
}

func genIf(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	t := typedecode.Block(load.Varint7())

	pushBranchTarget(f, t, false) // end
	var afterThen link.L

	cond := popOperand(f)
	if cond.Type != abi.I32 {
		panic(fmt.Errorf("if condition has wrong type: %s", cond.Type))
	}

	opSaveTemporaryOperands(f)
	opInitVars(f)
	opStoreVars(f, false)
	opBranchIf(f, cond, false, &afterThen)

	thenDeadend, haveElse := genThenOps(f, load)

	if !haveElse && t != abi.Void {
		panic(errors.New("if without else has a value type"))
	}

	if !thenDeadend {
		if t != abi.Void {
			value := popOperand(f)
			if value.Type != t {
				panic(fmt.Errorf("%s value operand has type %s, but target expects %s", op, value.Type, t))
			}
			opMove(f, reg.Result, value, false)
		}

		if haveElse {
			opSaveTemporaryOperands(f)
			opStoreVars(f, true)
			opBranch(f, &getBranchTarget(f, 0).Label) // end
		}
	}

	opLabel(f, &afterThen)
	isa.UpdateBranches(f.Text.Bytes(), &afterThen)

	if haveElse {
		deadend = genOps(f, load)

		if t != abi.Void && !deadend {
			value := popOperand(f)
			if value.Type != t {
				panic(fmt.Errorf("%s value operand has type %s, but target expects %s", op, value.Type, t))
			}
			opMove(f, reg.Result, value, false)
		}
	}

	end := popBranchTarget(f)
	if end.Live() { // includes thenDeadend information
		deadend = false
	}
	if !deadend {
		opLabel(f, end)
		isa.UpdateBranches(f.Text.Bytes(), end)
	}

	if t != abi.Void {
		pushResultRegOperand(f, t)
	}
	return
}

func genLoop(f *gen.Func, load loader.L, op Opcode, info opInfo) (deadend bool) {
	encodedBlockType := load.Varint7()

	opSaveTemporaryOperands(f)
	opInitVars(f)
	opStoreVars(f, false)

	pushBranchTarget(f, abi.Void, false) // begin
	opLabel(f, &getBranchTarget(f, 0).Label)

	savedMinBlockOperand := f.MinBlockOperand
	f.MinBlockOperand = len(f.Operands)

	deadend = genOps(f, load)

	if deadend {
		for len(f.Operands) > f.MinBlockOperand {
			x := popOperand(f)
			debugf("discarding operand at end of %s due to deadend: %s", op, x)
			discard(f, x)
		}
	} else {
		need := f.MinBlockOperand
		if encodedBlockType != 0 {
			need++ // result remains on stack
		}
		if len(f.Operands) > need { // let the next guy deal with missing operands
			panic(fmt.Errorf("operands remain on stack after %s", op))
		}
	}

	f.MinBlockOperand = savedMinBlockOperand

	begin := popBranchTarget(f)
	isa.UpdateBranches(f.Text.Bytes(), begin)
	return
}

func opLabel(f *gen.Func, l *link.L) {
	opSaveTemporaryOperands(f)
	opStoreVars(f, true)
	l.Addr = f.Text.Addr

	debugf("label")
}

func opBranch(f *gen.Func, l *link.L) {
	retAddr := asm.Branch(f.M, l.Addr)
	if l.Addr == 0 {
		l.AddSite(retAddr)
	}
}

func opBranchIf(f *gen.Func, x val.Operand, yes bool, l *link.L) {
	x = effectiveOperand(f, x)
	retAddrs := asm.BranchIf(f, x, yes, l.Addr)
	if l.Addr == 0 {
		l.AddSites(retAddrs)
	}
}

func opBranchIfOutOfBounds(f *gen.Func, indexReg reg.R, upperBound int32, l *link.L) {
	site := asm.BranchIfOutOfBounds(f.M, indexReg, upperBound, l.Addr)
	if l.Addr == 0 {
		l.AddSite(site)
	}
}
