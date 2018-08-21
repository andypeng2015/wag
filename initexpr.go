// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wag

import (
	"fmt"

	"github.com/tsavola/wag/internal/codegen"
	"github.com/tsavola/wag/internal/loader"
	"github.com/tsavola/wag/wasm"
)

func readInitExpr(m *Module, load loader.L) (valueBits uint64, t wasm.Type) {
	op := codegen.Opcode(load.Byte())

	switch op {
	case codegen.OpcodeI32Const:
		valueBits = uint64(int64(load.Varint32()))
		t = wasm.I32

	case codegen.OpcodeI64Const:
		valueBits = uint64(load.Varint64())
		t = wasm.I64

	case codegen.OpcodeF32Const:
		valueBits = uint64(load.Uint32())
		t = wasm.F32

	case codegen.OpcodeF64Const:
		valueBits = load.Uint64()
		t = wasm.F64

	case codegen.OpcodeGetGlobal:
		i := load.Varuint32()
		if i >= uint32(m.NumImportGlobals) {
			panic(fmt.Errorf("import global index out of bounds in initializer expression: %d", i))
		}
		g := m.Globals[i]
		valueBits = g.Init
		t = g.Type

	default:
		panic(fmt.Errorf("unsupported operation in initializer expression: %s", op))
	}

	if op := codegen.Opcode(load.Byte()); op != codegen.OpcodeEnd {
		panic(fmt.Errorf("unexpected operation in initializer expression when expecting end: %s", op))
	}

	return
}

func readOffsetInitExpr(m *Module, load loader.L) uint32 {
	offset, t := readInitExpr(m, load)
	if t != wasm.I32 {
		panic(fmt.Errorf("offset initializer expression has invalid type: %s", t))
	}
	return uint32(int32(int64(offset)))
}
