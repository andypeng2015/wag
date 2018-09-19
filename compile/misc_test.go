// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/tsavola/wag/internal/test/runner"
	"github.com/tsavola/wag/object/debug/dump"
	"github.com/tsavola/wag/static"
)

func TestCallWithDuplicatedOperand(t *testing.T) {
	misc(t, "../testdata/call-with-duplicated-operand.wast", "32744 32 32\n")
}

func misc(t *testing.T, filename, expectOutput string) {
	const (
		maxTextSize   = 65536
		maxRODataSize = 4096
		stackSize     = 4096

		dumpText = false
	)

	wasmData, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	wasmReadCloser := wast2wasm(wasmData, false)
	defer wasmReadCloser.Close()
	wasm := bufio.NewReader(wasmReadCloser)

	p, err := runner.NewProgram(maxTextSize, maxRODataSize)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	var mod = &Module{}
	mod.loadInitialSections(wasm)
	mod.defineImports(runner.Resolver)

	var code = &CodeConfig{
		Text:         static.Buf(p.Text),
		ROData:       static.Buf(p.ROData),
		RODataAddr:   p.RODataAddr(),
		ObjectMapper: &p.ObjInfo,
	}
	loadCodeSection(code, wasm, mod)

	var data = &DataConfig{}
	loadDataSection(data, wasm, mod)

	p.Seal()
	p.SetData(data.GlobalsMemory.Bytes(), mod.GlobalsSize())
	minMemorySize, maxMemorySize := mod.MemoryLimits()

	if dumpText && testing.Verbose() {
		dump.Text(os.Stdout, code.Text.Bytes(), p.TextAddr(), p.RODataAddr(), p.ObjInfo.FuncAddrs, nil)
	}

	var printBuf bytes.Buffer

	r, err := p.NewRunner(minMemorySize, maxMemorySize, stackSize)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	_, err = r.Run(0, mod.Sigs(), &printBuf)
	if err != nil {
		t.Fatal(err)
	}

	output := string(printBuf.Bytes())
	t.Logf("print output:\n%s", output)
	if output != expectOutput {
		t.Fail()
	}
}
