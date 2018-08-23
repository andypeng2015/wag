// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

//go:generate go run ../../internal/cmd/syscalls/generate.go

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/buffer"
	"github.com/tsavola/wag/compile"
)

var (
	verbose = false
)

var importFuncs = make(map[string]uint64)

type env struct{}

func (*env) ImportFunction(module, field string, sig abi.FunctionType) (variadic bool, absAddr uint64, err error) {
	if verbose {
		log.Printf("import %s%s", field, sig)
	}

	if module != "env" {
		err = fmt.Errorf("imported function's module is unknown: %s %s", module, field)
		return
	}

	absAddr = importFuncs[field]
	if absAddr == 0 {
		err = fmt.Errorf("imported function not supported: %s %s", module, field)
		return
	}

	return
}

func (*env) ImportGlobal(module, field string, t abi.Type) (valueBits uint64, err error) {
	err = fmt.Errorf("imported global not supported: %s %s", module, field)
	return
}

func makeMem(size int, extraProt, extraFlags int) (mem []byte, err error) {
	if size > 0 {
		mem, err = syscall.Mmap(-1, 0, size, syscall.PROT_READ|syscall.PROT_WRITE|extraProt, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS|extraFlags)
	}
	return
}

func memAddr(mem []byte) uintptr {
	return (*reflect.SliceHeader)(unsafe.Pointer(&mem)).Data
}

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] wasmfile\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	var (
		textSize    = 128 * 1024 * 1024
		roDataSize  = 4 * 1024 * 1024
		stackSize   = 64 * 1024
		entrySymbol = "main"
	)

	flag.BoolVar(&verbose, "v", verbose, "verbose logging")
	flag.IntVar(&textSize, "textsize", textSize, "maximum program text size")
	flag.IntVar(&roDataSize, "rodatasize", roDataSize, "maximum read-only data size")
	flag.IntVar(&stackSize, "stacksize", stackSize, "call stack size")
	flag.StringVar(&entrySymbol, "func", entrySymbol, "function to run")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	filename := flag.Arg(0)

	prog, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	progReader := bytes.NewReader(prog)

	roDataMem, err := makeMem(roDataSize, 0, syscall.MAP_32BIT)
	if err != nil {
		log.Fatal(err)
	}
	roDataAddr := memAddr(roDataMem)
	roDataBuf := buffer.NewFixed(roDataMem[:0])

	textMem, err := makeMem(textSize, syscall.PROT_EXEC, 0)
	if err != nil {
		log.Fatal(err)
	}
	textAddr := memAddr(textMem)
	textBuf := buffer.NewFixed(textMem[:0])

	pageSize := os.Getpagesize()
	pageMask := pageSize - 1

	m := compile.Module{
		EntrySymbol:     entrySymbol,
		EntryArgs:       make([]uint64, 2),
		MemoryAlignment: pageSize,
	}

	err = m.LoadPreliminarySections(progReader, &env{})
	if err != nil {
		log.Fatal(err)
	}

	err = m.LoadCodeSection(progReader, textBuf, roDataBuf, int32(roDataAddr), nil)
	if err != nil {
		log.Fatal(err)
	}

	memoryOffset := (m.GlobalsSize() + pageMask) &^ pageMask
	initMemorySize, growMemorySize := m.MemoryLimits()
	globalsMemoryMem, err := makeMem(memoryOffset+int(growMemorySize), 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	dataBuf := buffer.NewFixed(globalsMemoryMem[:0])
	memoryAddr := memAddr(globalsMemoryMem) + uintptr(memoryOffset)
	initMemoryEnd := memoryAddr + uintptr(initMemorySize)
	growMemoryEnd := memoryAddr + uintptr(growMemorySize)

	err = m.LoadDataSection(progReader, dataBuf)
	if err != nil {
		log.Fatal(err)
	}

	if err := syscall.Mprotect(roDataMem, syscall.PROT_READ); err != nil {
		log.Fatal(err)
	}

	if err := syscall.Mprotect(textMem, syscall.PROT_EXEC); err != nil {
		log.Fatal(err)
	}

	stackMem, err := makeMem(stackSize, 0, syscall.MAP_STACK)
	if err != nil {
		log.Fatal(err)
	}
	stackAddr := memAddr(stackMem)
	stackEnd := stackAddr + uintptr(stackSize)

	exec(textAddr, stackAddr, memoryAddr, initMemoryEnd, growMemoryEnd, stackEnd)
}
