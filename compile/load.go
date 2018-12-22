// Copyright (c) 2015 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compile implements a WebAssembly compiler.
package compile

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"

	"github.com/tsavola/wag/buffer"
	"github.com/tsavola/wag/compile/event"
	"github.com/tsavola/wag/internal/code"
	"github.com/tsavola/wag/internal/data"
	"github.com/tsavola/wag/internal/datalayout"
	"github.com/tsavola/wag/internal/errorpanic"
	"github.com/tsavola/wag/internal/gen/codegen"
	"github.com/tsavola/wag/internal/initexpr"
	"github.com/tsavola/wag/internal/loader"
	"github.com/tsavola/wag/internal/module"
	"github.com/tsavola/wag/internal/obj"
	"github.com/tsavola/wag/internal/reader"
	"github.com/tsavola/wag/internal/section"
	"github.com/tsavola/wag/internal/typedecode"
	"github.com/tsavola/wag/wa"
)

const (
	DefaultMaxTextSize = 0x7fff0000 // below 2 GB to mitigate address calculation bugs

	defaultTextBufferSize   = wa.PageSize // arbitrary
	defaultMemoryBufferSize = wa.PageSize // arbitrary
)

var emptyCodeSection = []byte{
	1, // payload length
	0, // function count
}

// Reader is a subset of bufio.Reader, bytes.Buffer and bytes.Reader.
type Reader = reader.R

type CodeBuffer = code.Buffer
type DataBuffer = data.Buffer

const (
	maxStringLen          = 255   // TODO
	maxTableLimit         = 32768 // TODO
	maxInitialMemoryLimit = 16384 // TODO
	maxMaximumMemoryLimit = math.MaxInt32 >> wa.PageBits
	maxGlobals            = 4096/obj.Word - 2 // (trap handler + memory limit)
	maxExports            = 256
	maxElements           = 32768
)

func readResizableLimits(load loader.L, maxInitial, maxMaximum uint32, scale int) module.ResizableLimits {
	maximumFieldIsPresent := load.Varuint1()

	initial := load.Varuint32()
	if initial > maxInitial {
		panic(module.Errorf("initial memory size is too large: %d", initial))
	}

	maximum := maxMaximum

	if maximumFieldIsPresent {
		maximum = load.Varuint32()
		if maximum > maxMaximum {
			maximum = maxMaximum
		}
		if maximum < initial {
			panic(module.Errorf("maximum memory size %d is smaller than initial memory size %d", maximum, initial))
		}
	}

	return module.ResizableLimits{
		Initial: int(initial) * scale,
		Maximum: int(maximum) * scale,
		Defined: true,
	}
}

func readTable(m *module.M, load loader.L) {
	if m.TableLimitValues.Defined {
		panic(module.Error("multiple tables not supported"))
	}

	if elementType := load.Varint7(); elementType != -0x10 {
		panic(module.Errorf("unsupported table element type: %d", elementType))
	}

	m.TableLimitValues = readResizableLimits(load, maxTableLimit, maxTableLimit, 1)
}

func readMemory(m *module.M, load loader.L) {
	if m.MemoryLimitValues.Defined {
		panic(module.Error("multiple memories not supported"))
	}

	m.MemoryLimitValues = readResizableLimits(load, maxInitialMemoryLimit, maxMaximumMemoryLimit, wa.PageSize)
}

// Config for loading WebAssembly module sections.
type Config struct {
	// SectionMapper is invoked for every section (standard or custom), just
	// after the section id byte.  It must read and return the payload length
	// (varuint32), but not the payload itself.
	SectionMapper func(sectionId byte, r Reader) (payloadLen uint32, err error)

	// CustomSectionLoader is invoked for every custom section.  It must read
	// exactly payloadLen bytes, or return an error.  SectionMapper (if
	// configured) has been invoked just before it.
	CustomSectionLoader func(r Reader, payloadLen uint32) error
}

// ModuleConfig for a single compiler invocation.
type ModuleConfig struct {
	Config
}

// Module contains a WebAssembly module specification without code or data.
type Module struct {
	m module.M
}

// LoadInitialSections reads module header and all sections preceding code and
// data.
func LoadInitialSections(config *ModuleConfig, r Reader) (m *Module, err error) {
	defer func() {
		err = errorpanic.Handle(recover())
	}()

	m = loadInitialSections(config, r)
	return
}

func loadInitialSections(config *ModuleConfig, r Reader) (m *Module) {
	if config == nil {
		config = new(ModuleConfig)
	}

	m = new(Module)
	load := loader.L{R: r}

	var header module.Header
	if err := binary.Read(load.R, binary.LittleEndian, &header); err != nil {
		panic(err)
	}
	if header.MagicNumber != module.MagicNumber {
		panic(module.Error("not a WebAssembly module"))
	}
	if header.Version != module.Version {
		panic(module.Errorf("unsupported module version: %d", header.Version))
	}

	var seenId module.SectionId

	for {
		sectionId, err := load.R.ReadByte()
		if err != nil {
			if err == io.EOF {
				return
			}
			panic(err)
		}

		id := module.SectionId(sectionId)

		if id != module.SectionCustom {
			if id <= seenId {
				panic(module.Errorf("section 0x%x follows section 0x%x", id, seenId))
			}
			seenId = id
		}

		if id >= module.NumMetaSections {
			load.R.UnreadByte()
			if id >= module.NumSections {
				panic(module.Errorf("custom section id: 0x%x", id))
			}
			return
		}

		var payloadLen uint32

		if config.SectionMapper != nil {
			payloadLen, err = config.SectionMapper(sectionId, r)
			if err != nil {
				panic(err)
			}
		} else {
			payloadLen = load.Varuint32()
		}

		metaSectionLoaders[id](m, config, payloadLen, load)
	}
}

var metaSectionLoaders = [module.NumMetaSections]func(*Module, *ModuleConfig, uint32, loader.L){
	module.SectionCustom: func(m *Module, config *ModuleConfig, payloadLen uint32, load loader.L) {
		var err error
		if config.CustomSectionLoader != nil {
			err = config.CustomSectionLoader(load.R, payloadLen)
		} else {
			_, err = io.CopyN(ioutil.Discard, load.R, int64(payloadLen))
		}
		if err != nil {
			panic(err)
		}
	},

	module.SectionType: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for i := range load.Count(module.MaxTypes, "type") {
			if form := load.Varint7(); form != -0x20 {
				panic(module.Errorf("unsupported function type form: %d", form))
			}

			var sig wa.FuncType

			paramCount := load.Varuint32()
			if paramCount > module.MaxFuncParams {
				panic(module.Errorf("function type #%d has too many parameters: %d", i, paramCount))
			}

			sig.Params = make([]wa.Type, paramCount)
			for j := range sig.Params {
				sig.Params[j] = typedecode.Value(load.Varint7())
			}

			if returnCount1 := load.Varuint1(); returnCount1 {
				sig.Result = typedecode.Value(load.Varint7())
			}

			m.m.Types = append(m.m.Types, sig)
		}
	},

	module.SectionImport: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for i := range load.Count(module.MaxImports, "import") {
			moduleLen := load.Varuint32()
			if moduleLen > maxStringLen {
				panic(module.Errorf("module string is too long in import #%d", i))
			}

			moduleStr := string(load.Bytes(moduleLen))

			fieldLen := load.Varuint32()
			if fieldLen > maxStringLen {
				panic(module.Errorf("field string is too long in import #%d", i))
			}

			fieldStr := string(load.Bytes(fieldLen))

			kind := module.ExternalKind(load.Byte())

			switch kind {
			case module.ExternalKindFunction:
				sigIndex := load.Varuint32()
				if sigIndex >= uint32(len(m.m.Types)) {
					panic(module.Errorf("function type index out of bounds in import #%d: 0x%x", i, sigIndex))
				}

				m.m.Funcs = append(m.m.Funcs, sigIndex)

				m.m.ImportFuncs = append(m.m.ImportFuncs, module.ImportFunc{
					Import: module.Import{
						Module: moduleStr,
						Field:  fieldStr,
					},
				})

			case module.ExternalKindTable:
				readTable(&m.m, load)

			case module.ExternalKindMemory:
				readMemory(&m.m, load)

			case module.ExternalKindGlobal:
				if len(m.m.Globals) >= maxGlobals {
					panic(module.Error("too many imported globals"))
				}

				t := typedecode.Value(load.Varint7())

				if mutable := load.Varuint1(); mutable {
					panic(module.Errorf("unsupported mutable global in import #%d", i))
				}

				m.m.Globals = append(m.m.Globals, module.Global{
					Type: t,
				})

				m.m.ImportGlobals = append(m.m.ImportGlobals, module.Import{
					Module: moduleStr,
					Field:  fieldStr,
				})

			default:
				panic(module.Errorf("import kind not supported: %s", kind))
			}
		}
	},

	module.SectionFunction: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for range load.Count(module.MaxFunctions, "function") {
			sigIndex := load.Varuint32()
			if sigIndex >= uint32(len(m.m.Types)) {
				panic(module.Errorf("function type index out of bounds: %d", sigIndex))
			}

			m.m.Funcs = append(m.m.Funcs, sigIndex)
		}
	},

	module.SectionTable: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for range load.Count(1, "table") {
			readTable(&m.m, load)
		}
	},

	module.SectionMemory: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for range load.Count(1, "memory") {
			readMemory(&m.m, load)
		}
	},

	module.SectionGlobal: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for range load.Count(maxGlobals, "global") {
			t := typedecode.Value(load.Varint7())
			mutable := load.Varuint1()
			init, _ := initexpr.Read(&m.m, load)

			m.m.Globals = append(m.m.Globals, module.Global{
				Type:    t,
				Mutable: mutable,
				Init:    init,
			})
		}
	},

	module.SectionExport: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		m.m.ExportFuncs = make(map[string]uint32)

		for i := range load.Count(maxExports, "export") {
			fieldLen := load.Varuint32()
			if fieldLen > maxStringLen {
				panic(module.Errorf("field string is too long in export #%d", i))
			}

			fieldStr := load.Bytes(fieldLen)
			kind := module.ExternalKind(load.Byte())
			index := load.Varuint32()

			switch kind {
			case module.ExternalKindFunction:
				if index >= uint32(len(m.m.Funcs)) {
					panic(module.Errorf("export function index out of bounds: %d", index))
				}
				m.m.ExportFuncs[string(fieldStr)] = index

			case module.ExternalKindTable, module.ExternalKindMemory, module.ExternalKindGlobal:

			default:
				panic(module.Errorf("custom export kind: %s", kind))
			}
		}
	},

	module.SectionStart: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		index := load.Varuint32()
		if index >= uint32(len(m.m.Funcs)) {
			panic(module.Errorf("start function index out of bounds: %d", index))
		}

		sigIndex := m.m.Funcs[index]
		sig := m.m.Types[sigIndex]
		if len(sig.Params) > 0 || sig.Result != wa.Void {
			panic(module.Errorf("invalid start function signature: %s", sig))
		}

		m.m.StartIndex = index
		m.m.StartDefined = true
	},

	module.SectionElement: func(m *Module, _ *ModuleConfig, _ uint32, load loader.L) {
		for i := range load.Count(maxElements, "element") {
			if index := load.Varuint32(); index != 0 {
				panic(module.Errorf("unsupported table index: %d", index))
			}

			offset := initexpr.ReadOffset(&m.m, load)

			numElem := load.Varuint32()

			needSize := uint64(offset) + uint64(numElem)
			if needSize > uint64(m.m.TableLimitValues.Initial) {
				panic(module.Errorf("table segment #%d exceeds initial table size", i))
			}

			oldSize := len(m.m.TableFuncs)
			if needSize > uint64(oldSize) {
				buf := make([]uint32, needSize)
				copy(buf, m.m.TableFuncs)
				for i := oldSize; i < int(offset); i++ {
					buf[i] = math.MaxInt32 // invalid function index
				}
				m.m.TableFuncs = buf
			}

			for j := int(offset); j < int(needSize); j++ {
				elem := load.Varuint32()
				if elem >= uint32(len(m.m.Funcs)) {
					panic(module.Errorf("table element index out of bounds: %d", elem))
				}

				m.m.TableFuncs[j] = elem
			}
		}
	},
}

func (m *Module) Types() []wa.FuncType      { return m.m.Types }
func (m *Module) FuncTypeIndexes() []uint32 { return m.m.Funcs }

func (m *Module) FuncTypes() []wa.FuncType {
	sigs := make([]wa.FuncType, len(m.m.Funcs))
	for i, sigIndex := range m.m.Funcs {
		sigs[i] = m.m.Types[sigIndex]
	}
	return sigs
}

func (m *Module) InitialMemorySize() int { return m.m.MemoryLimitValues.Initial }
func (m *Module) MemorySizeLimit() int   { return m.m.MemoryLimitValues.Maximum }

func (m *Module) NumImportFuncs() int   { return len(m.m.ImportFuncs) }
func (m *Module) NumImportGlobals() int { return len(m.m.ImportGlobals) }

func (m *Module) ImportFunc(i int) (module, field string, sig wa.FuncType) {
	imp := m.m.ImportFuncs[i]
	module = imp.Module
	field = imp.Field

	sigIndex := m.m.Funcs[i]
	sig = m.m.Types[sigIndex]
	return
}

func (m *Module) ImportGlobal(i int) (module, field string, t wa.Type) {
	imp := m.m.ImportGlobals[i]
	module = imp.Module
	field = imp.Field

	t = m.m.Globals[i].Type
	return
}

func (m *Module) SetImportFunc(i int, vecIndex int)  { m.m.ImportFuncs[i].VecIndex = vecIndex }
func (m *Module) SetImportGlobal(i int, init uint64) { m.m.Globals[i].Init = init }

func (m *Module) GlobalsSize() int {
	size := len(m.m.Globals) * obj.Word
	mask := datalayout.MinAlignment - 1 // Round up so that linear memory will
	return (size + mask) &^ mask        // have at least minimum alignment.
}

func (m *Module) ExportFuncs() map[string]uint32 { return m.m.ExportFuncs }

func (m *Module) ExportFunc(field string) (funcIndex uint32, sig wa.FuncType, found bool) {
	funcIndex, found = m.m.ExportFuncs[field]
	if found {
		sigIndex := m.m.Funcs[funcIndex]
		sig = m.m.Types[sigIndex]
	}
	return
}

// CodeConfig for a single compiler invocation.  Either MaxTextSize or Text
// should be specified, but not both.
type CodeConfig struct {
	MaxTextSize  int        // Effective if Text is nil; defaults to DefaultMaxTextSize.
	Text         CodeBuffer // Initialized with default implementation if nil.
	Mapper       ObjectMapper
	EventHandler func(event.Event)
	LastInitFunc uint32
	Config
}

// LoadCodeSection reads a WebAssembly module's code section and generates
// machine code.
//
// If CodeBuffer panics with an error, it will be returned by this function.
func LoadCodeSection(config *CodeConfig, r Reader, mod *Module) (err error) {
	defer func() {
		err = errorpanic.Handle(recover())
	}()

	loadCodeSection(config, r, mod)
	return
}

func loadCodeSection(config *CodeConfig, r Reader, mod *Module) {
	load := loader.L{R: r}

	switch id := section.Find(module.SectionCode, load, config.SectionMapper, config.CustomSectionLoader); id {
	case module.SectionData, 0:
		// no code section

		load = loader.L{R: bytes.NewReader(emptyCodeSection)}
		fallthrough

	case module.SectionCode:
		var payloadLen uint32
		var err error

		if config.SectionMapper != nil {
			payloadLen, err = config.SectionMapper(byte(id), load.R)
			if err != nil {
				panic(err)
			}
		} else {
			payloadLen = load.Varuint32()
		}

		if config.Text == nil {
			if config.MaxTextSize == 0 {
				config.MaxTextSize = DefaultMaxTextSize
			}

			alloc := defaultTextBufferSize
			if alloc > config.MaxTextSize {
				alloc = config.MaxTextSize
			}
			if guess := 512 + uint64(payloadLen)*8; guess < uint64(alloc) { // citation needed
				alloc = int(guess)
			}

			config.Text = buffer.NewLimited(make([]byte, 0, alloc), config.MaxTextSize)
		}

		mapper := config.Mapper
		if mapper == nil {
			mapper = dummyMap{}
		}

		codegen.GenProgram(config.Text, mapper, load, &mod.m, config.EventHandler, int(config.LastInitFunc)+1)

	default:
		panic(module.Errorf("unexpected section id: 0x%x (looking for code section)", id))
	}
}

// DataConfig for a single compiler invocation.
type DataConfig struct {
	GlobalsMemory   DataBuffer // Initialized with default implementation if nil.
	MemoryAlignment int        // Initialized with minimal value if zero.
	Config
}

// LoadDataSection reads a WebAssembly module's data section and generates
// initial contents of mutable program state (globals and linear memory).
//
// If DataBuffer panics with an error, it will be returned by this function.
func LoadDataSection(config *DataConfig, r Reader, mod *Module) (err error) {
	defer func() {
		err = errorpanic.Handle(recover())
	}()

	loadDataSection(config, r, mod)
	return
}

func loadDataSection(config *DataConfig, r Reader, mod *Module) {
	if config.MemoryAlignment == 0 {
		config.MemoryAlignment = datalayout.MinAlignment
	}
	memoryOffset := datalayout.MemoryOffset(&mod.m, config.MemoryAlignment)

	load := loader.L{R: r}

	switch id := section.Find(module.SectionData, load, config.SectionMapper, config.CustomSectionLoader); id {
	case module.SectionData:
		var payloadLen uint32
		var err error

		if config.SectionMapper != nil {
			payloadLen, err = config.SectionMapper(byte(id), load.R)
			if err != nil {
				panic(err)
			}
		} else {
			payloadLen = load.Varuint32()
		}

		if config.GlobalsMemory == nil {
			memAlloc := defaultMemoryBufferSize
			if payloadLen < uint32(memAlloc) {
				memAlloc = int(payloadLen) // hope for dense packing
			}

			limit := memoryOffset + mod.InitialMemorySize()
			alloc := memoryOffset + memAlloc
			if alloc > limit {
				alloc = limit
			}

			config.GlobalsMemory = buffer.NewDynamicHint(make([]byte, 0, alloc), limit)
		}

		datalayout.CopyGlobalsAlign(config.GlobalsMemory, &mod.m, memoryOffset)
		datalayout.ReadMemory(config.GlobalsMemory, load, &mod.m)

	case 0:
		// no data section

		if config.GlobalsMemory == nil {
			config.GlobalsMemory = buffer.NewStatic(make([]byte, 0, memoryOffset), memoryOffset)
		}

		datalayout.CopyGlobalsAlign(config.GlobalsMemory, &mod.m, memoryOffset)

	default:
		panic(module.Errorf("unexpected section id: 0x%x (looking for data section)", id))
	}
}

// ValidateDataSection reads a WebAssembly module's data section.
func ValidateDataSection(config *Config, r Reader, mod *Module) (err error) {
	defer func() {
		err = errorpanic.Handle(recover())
	}()

	validateDataSection(config, r, mod)
	return
}

func validateDataSection(config *Config, r Reader, mod *Module) {
	if config == nil {
		config = new(Config)
	}

	load := loader.L{R: r}

	switch id := section.Find(module.SectionData, load, config.SectionMapper, config.CustomSectionLoader); id {
	case module.SectionData:
		if config.SectionMapper != nil {
			_, err := config.SectionMapper(byte(id), load.R)
			if err != nil {
				panic(err)
			}
		} else {
			load.Varuint32()
		}

		datalayout.ValidateMemory(load, &mod.m)

	case 0:
		// no data section

	default:
		panic(module.Errorf("unexpected section id: 0x%x (looking for data section)", id))
	}
}

// LoadCustomSections reads a WebAssembly module's extension sections which
// follow standard sections.
func LoadCustomSections(config *Config, r Reader) (err error) {
	defer func() {
		err = errorpanic.Handle(recover())
	}()

	loadCustomSections(config, r)
	return
}

func loadCustomSections(config *Config, r Reader) {
	if config == nil {
		config = new(Config)
	}

	load := loader.L{R: r}

	if id := section.Find(0, load, config.SectionMapper, config.CustomSectionLoader); id != 0 {
		panic(module.Errorf("unexpected section id: 0x%x (after all standard sections)", id))
	}
}
