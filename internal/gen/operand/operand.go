// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package operand

import (
	"fmt"

	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/internal/gen/condition"
	"github.com/tsavola/wag/internal/gen/debug"
	"github.com/tsavola/wag/internal/gen/reg"
	"github.com/tsavola/wag/internal/gen/storage"
	"github.com/tsavola/wag/internal/isa/reglayout"
)

const (
	payloadZeroExt = uint64(1 << reglayout.Radix)
	payloadDebug   = uint64(0x0c7a0c707e7a7e55) // "octa octo teta tess"
)

type O struct {
	Storage storage.Storage
	Type    abi.Type

	payload uint64
}

func Placeholder(t abi.Type) (o O) {
	o.Storage = storage.Imm // Least dangerous
	o.Type = t
	if debug.Enabled {
		o.payload = payloadDebug
	}
	return
}

func Imm(t abi.Type, value uint64) O {
	return O{storage.Imm, t, value}
}

func Reg(t abi.Type, r reg.R, zeroExt bool) O {
	value := uint64(byte(r))
	if zeroExt {
		value |= payloadZeroExt
	}
	return O{storage.Reg, t, value}
}

func Flags(cond condition.C) O {
	return O{storage.Flags, abi.I32, uint64(cond)}
}

func (o *O) SetPlaceholder() {
	o.Storage = storage.Imm // Least dangerous
	if debug.Enabled {
		o.payload = payloadDebug
	}
}

func (o *O) SetReg(r reg.R, zeroExt bool) {
	o.Storage = storage.Reg
	o.payload = uint64(r)
	if zeroExt {
		o.payload |= payloadZeroExt
	}
}

func (o *O) SetStack() {
	o.Storage = storage.Stack
}

func (o O) ImmValue() int64 {
	if o.Type.Size() == abi.Size32 {
		return int64(int32(uint32(o.payload)))
	} else {
		return int64(o.payload)
	}
}

func (o O) ImmValue8() int8        { return int8(uint8(o.payload)) }
func (o O) Reg() reg.R             { return reg.R(byte(o.payload)) }
func (o O) RegZeroExt() bool       { return o.payload >= payloadZeroExt }
func (o O) FlagsCond() condition.C { return condition.C(int(o.payload)) }

func (o O) String() string {
	if debug.Enabled && o.payload == payloadDebug {
		if o.Storage == storage.Imm {
			return fmt.Sprintf("%s placeholder", o.Type)
		} else {
			return fmt.Sprintf("magic %s value 0x%08x surprisingly in %s storage", o.Type, o.payload, o.Storage)
		}
	}

	switch o.Storage {
	case storage.Imm:
		if o.Type.Category() == abi.Int {
			return fmt.Sprintf("constant %s 0x%x", o.Type, o.payload)
		} else {
			return fmt.Sprintf("constant %s bits 0x%x", o.Type, o.payload)
		}

	case storage.Stack:
		return fmt.Sprintf("saved %s in stack", o.Type)

	case storage.Reg:
		if o.Reg() == reg.Result {
			return fmt.Sprintf("volatile %s in result register", o.Type)
		} else {
			return fmt.Sprintf("allocated %s in %s", o.Type, o.Reg())
		}

	case storage.Flags:
		return fmt.Sprintf("volatile %s comparison result", o.FlagsCond())

	default:
		return "<invalid operand>"
	}
}
