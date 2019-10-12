// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package wag provides a high-level WebAssembly compiler API.

See the Compile function's source code for an example of how to use the
low-level compiler APIs (implemented in subpackages).


Errors

ModuleError type is accessible via errors subpackage.  Such errors may be
returned by compilation and other parsing functions.  Other types of errors
indicate either a read error or an internal compiler error.  Unexpected EOF is
a ModuleError which wraps io.ErrUnexpectedEOF.

Default buffer implementations use the buffer.ErrSizeLimit error to indicate
that generated code or data doesn't fit in a target buffer.  It is a
ModuleError (the module didn't conform to size constraints).

*/
package wag

import (
	"fmt"

	"github.com/tsavola/wag/binding"
	"github.com/tsavola/wag/compile"
	"github.com/tsavola/wag/object/debug"
	"github.com/tsavola/wag/object/stack"
	"github.com/tsavola/wag/section"
	"github.com/tsavola/wag/wa"
)

func CompileLibrary(r compile.Reader, imports binding.LibraryImportResolver) (lib compile.Library, err error) {
	mod, err := compile.LoadInitialSections(nil, r)
	if err != nil {
		return
	}

	lib, err = mod.AsLibrary()
	if err != nil {
		return
	}

	err = binding.BindLibraryImports(&lib, imports)
	if err != nil {
		return
	}

	err = lib.LoadSections(r)
	if err != nil {
		return
	}

	return
}

// EntryPolicy validates an entry function's signature while looking it up from
// a module's exported functions.
type EntryPolicy func(m compile.Module, symbol string) (globalIndex uint32, err error)

// Config for a single compiler invocation.  Zero values are replaced with
// effective defaults during compilation.
type Config struct {
	ImportResolver  binding.ImportResolver // Imports are mapped to the library by default.
	Text            compile.CodeBuffer     // Defaults to dynamically sized buffer.
	GlobalsMemory   compile.DataBuffer     // Defaults to dynamically sized buffer.
	MemoryAlignment int                    // Defaults to minimal valid alignment.
	Entry           string                 // No entry function by default.
}

// Object code with debug information.  The fields are roughly in order of
// appearance during compilation.
//
// Executing the code requires a platform-specific mechanism.  It is not
// supported by this package, but see the compile subpackage for information.
type Object struct {
	FuncTypes         []wa.FuncType       // Signatures for debug output.
	InitialMemorySize int                 // Current memory allocation.
	MemorySizeLimit   int                 // Maximum valid value if not limited.
	Text              []byte              // Machine code and read-only data.
	debug.InsnMap                         // Stack unwinding and debug metadata.
	MemoryOffset      int                 // Threshold between globals and memory.
	GlobalsMemory     []byte              // Global values and memory contents.
	StackFrame        []byte              // Start and entry function addresses.
	Names             section.NameSection // Symbols for debug output.
}

// Compile a WebAssembly binary module into machine code.  The Object is
// constructed incrementally so that populated fields may be inspected on
// error.
//
// See the source code for examples of how to use the lower-level APIs.
func Compile(objectConfig *Config, r compile.Reader, lib compile.Library) (object *Object, err error) {
	if objectConfig == nil {
		objectConfig = new(Config)
	}

	object = new(Object)

	// In general, custom sections may appear at any position in the binary
	// module, so the custom section loader must be available at every step.
	// (WebAssembly specification says that the name section can appear only
	// after the data section, but wag's custom section handling is decoupled
	// from standard section handling; just accept it at any point.)

	var customSections = section.CustomLoaders{
		section.CustomName: object.Names.Load,
	}

	var loadingConfig = compile.Config{
		CustomSectionLoader: customSections.Load,
	}

	// Construct the Module object while reading the WebAssembly sections
	// preceding the actual program code.  (The Module object needs to be
	// available during compilation and when looking up entry functions, but
	// the program can later be executed without it.)

	var moduleConfig = &compile.ModuleConfig{
		Config: loadingConfig,
	}

	module, err := compile.LoadInitialSections(moduleConfig, r)
	object.FuncTypes = module.FuncTypes()
	object.InitialMemorySize = module.InitialMemorySize()
	object.MemorySizeLimit = module.MemorySizeLimit()
	if err != nil {
		return
	}

	// Fill in host function addresses and global variables' values.

	if objectConfig.ImportResolver == nil {
		objectConfig.ImportResolver = resolver{lib}
	}

	err = binding.BindImports(&module, objectConfig.ImportResolver)
	if err != nil {
		return
	}

	// Generate executable code and debug information while reading the
	// WebAssembly code section.  Text encodes the import function vector
	// indexes, but not the function addresses (the vector can be mapped
	// separately during execution).  It is also independent of entry function
	// choice and program state.

	var codeConfig = &compile.CodeConfig{
		Text:   objectConfig.Text,
		Mapper: &object.CallMap,
		Config: loadingConfig,
	}

	err = compile.LoadCodeSection(codeConfig, r, module, lib)
	objectConfig.Text = codeConfig.Text
	object.Text = codeConfig.Text.Bytes()
	if err != nil {
		return
	}

	// Generate initial linear memory contents while reading the WebAssembly
	// data section.  This step also copies the global variables' initial
	// values into the same buffer, just before the memory contents.
	// MemoryAlignment causes padding to be inserted before the globals.

	var dataConfig = &compile.DataConfig{
		GlobalsMemory:   objectConfig.GlobalsMemory,
		MemoryAlignment: objectConfig.MemoryAlignment,
		Config:          loadingConfig,
	}

	err = compile.LoadDataSection(dataConfig, r, module)
	objectConfig.GlobalsMemory = dataConfig.GlobalsMemory
	objectConfig.MemoryAlignment = dataConfig.MemoryAlignment
	object.MemoryOffset = alignSize(module.GlobalsSize(), dataConfig.MemoryAlignment)
	object.GlobalsMemory = dataConfig.GlobalsMemory.Bytes()
	if err != nil {
		return
	}

	// Find the optional start function, and the export function which will be
	// used as the optional entry point.  CallMap is used to look up their
	// addresses.

	var startAddr uint32

	if startIndex, defined := module.StartFunc(); defined {
		startAddr = object.FuncAddrs[startIndex]
	}

	var (
		entryIndex uint32
		entryAddr  uint32
	)

	if objectConfig.Entry != "" {
		entryIndex, err = binding.EntryFunc(module, objectConfig.Entry)
		if err != nil {
			return
		}

		entryAddr = object.FuncAddrs[entryIndex]
	}

	// Form a stack frame for the init routine which calls the functions.

	object.StackFrame = stack.InitFrame(startAddr, entryAddr)

	// Read the whole binary module to get the name section.

	err = compile.LoadCustomSections(&loadingConfig, r)
	if err != nil {
		return
	}

	return
}

// alignSize rounds up.
func alignSize(size, alignment int) int {
	return (size + (alignment - 1)) &^ (alignment - 1)
}

// resolver looks up program module's imports from the intermediate
// library module.
type resolver struct {
	lib compile.Library
}

func (r resolver) ResolveFunc(module, field string, sig wa.FuncType) (funcIndex uint32, err error) {
	funcIndex, actualSig, fieldFound := r.lib.ExportFunc(field)
	if module != "env" || !fieldFound {
		err = fmt.Errorf("unknown function imported: %q.%q", module, field)
		return
	}

	if !sig.Equal(actualSig) {
		err = fmt.Errorf("function %s.%s%s imported with wrong type: %s", module, field, actualSig, sig)
		return
	}

	return
}

func (r resolver) ResolveGlobal(module, field string, t wa.Type) (init uint64, err error) {
	// Globals are not supported by library.
	err = fmt.Errorf("unknown global imported: %q.%q", module, field)
	return
}
