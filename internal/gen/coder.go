// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"encoding/binary"

	"github.com/tsavola/wag/internal/regs"
	"github.com/tsavola/wag/internal/values"
	"github.com/tsavola/wag/trap"
	"github.com/tsavola/wag/wasm"
)

const (
	// Masks are 16-byte aligned for x86-64 SSE.
	ROMask7fAddr32 = iota * 16
	ROMask7fAddr64
	ROMask80Addr32
	ROMask80Addr64
	ROMask5f00Addr32 // 01011111000000000000000000000000
	ROMask43e0Addr64 // 0100001111100000000000000000000000000000000000000000000000000000
	ROTableAddr
)

const (
	WordSize     = 8              // stack entry size
	StackReserve = WordSize + 128 // trap/import call return address + red zone
)

type Buffer interface {
	Bytes() []byte
	Pos() int32
	Extend(n int) []byte
	PutByte(byte)
	PutBytes([]byte)
}

type Coder interface {
	Buffer

	MinMemorySize() int
	RODataAddr() int32
	TrapTrampolineAddr(id trap.Id) int32
	OpTrapCall(id trap.Id)

	Discard(values.Operand)
	Consumed(values.Operand)
	RegAllocated(wasm.Type, regs.R) bool
	FreeReg(wasm.Type, regs.R)
}

type RegCoder interface {
	Coder

	TryAllocReg(t wasm.Type) (reg regs.R, ok bool)
	AllocSpecificReg(t wasm.Type, reg regs.R)
}

type MaskBaseAddr int32

const (
	Mask7fBase    = MaskBaseAddr(ROMask7fAddr32)
	Mask80Base    = MaskBaseAddr(ROMask80Addr32)
	MaskTruncBase = MaskBaseAddr(ROMask5f00Addr32)
)

// MaskAddr calculates the absolute read-only data address for reading a mask
// for the given type size.  maskBaseAddr should be one of the Mask*Base
// constants.
func MaskAddr(roDataAddr int32, maskBaseAddr MaskBaseAddr, t wasm.Type) int32 {
	return roDataAddr + int32(maskBaseAddr) + int32((t.Size()&8)<<1)
}

func PutInt8(code Buffer, value int8) {
	code.PutByte(uint8(value))
}

func PutInt16(code Buffer, value int16) {
	binary.LittleEndian.PutUint16(code.Extend(2), uint16(value))
}

func PutInt32(code Buffer, value int32) {
	binary.LittleEndian.PutUint32(code.Extend(4), uint32(value))
}

func PutInt64(code Buffer, value int64) {
	binary.LittleEndian.PutUint64(code.Extend(8), uint64(value))
}
