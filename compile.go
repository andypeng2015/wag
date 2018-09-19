// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wag

import (
	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/compile"
	"github.com/tsavola/wag/object"
	"github.com/tsavola/wag/section"
	"github.com/tsavola/wag/wasm"
)

// Object code with stack map and debug symbols.  The text (machine code),
// read-only data, initial global values, and initial linear memory contents
// can be used to execute the program (the details are architecture-specific).
type Object struct {
	FuncSigs        []abi.Sig
	InitMemorySize  wasm.MemorySize
	GrowMemoryLimit wasm.MemorySize
	Text            []byte
	ROData          []byte
	object.CallMap
	MemoryOffset  int // Threshold between globals and memory.
	GlobalsMemory []byte
	Names         section.NameSection
}

// Config for a single compiler invocation.  Zero values are replaced with
// effective defaults during compilation.
type Config struct {
	EntrySymbol     string
	EntryArgs       []uint64
	Text            compile.CodeBuffer
	ROData          compile.DataBuffer
	RODataAddr      uintptr
	GlobalsMemory   compile.DataBuffer
	MemoryAlignment int
}

// Compile a WebAssembly binary module into machine code.
func Compile(config *Config, r compile.Reader, res compile.ImportResolver) (obj *Object, err error) {
	obj = new(Object)
	nameLoader := section.UnknownLoaders{"name": obj.Names.Load}.Load

	var mod = &compile.Module{
		EntrySymbol:          config.EntrySymbol,
		EntryArgs:            config.EntryArgs,
		UnknownSectionLoader: nameLoader,
	}

	err = mod.LoadInitialSections(r)
	if err != nil {
		return
	}

	err = mod.DefineImports(res)
	if err != nil {
		return
	}

	obj.FuncSigs = mod.FuncSigs()
	obj.InitMemorySize, obj.GrowMemoryLimit = mod.MemoryLimits()

	var code = &compile.CodeConfig{
		Text:                 config.Text,
		ROData:               config.ROData,
		RODataAddr:           config.RODataAddr,
		ObjectMapper:         &obj.CallMap,
		UnknownSectionLoader: nameLoader,
	}

	err = compile.LoadCodeSection(code, r, mod)
	if err != nil {
		return
	}

	config.Text = code.Text
	config.ROData = code.ROData
	obj.Text = code.Text.Bytes()
	obj.ROData = code.ROData.Bytes()

	var data = &compile.DataConfig{
		GlobalsMemory:        config.GlobalsMemory,
		MemoryAlignment:      config.MemoryAlignment,
		UnknownSectionLoader: nameLoader,
	}

	err = compile.LoadDataSection(data, r, mod)
	if err != nil {
		return
	}

	config.GlobalsMemory = data.GlobalsMemory
	config.MemoryAlignment = data.MemoryAlignment
	obj.MemoryOffset = (mod.GlobalsSize() + (data.MemoryAlignment - 1)) &^ (data.MemoryAlignment - 1)
	obj.GlobalsMemory = data.GlobalsMemory.Bytes()

	err = compile.LoadUnknownSections(nameLoader, r)
	return
}
