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

	"github.com/tsavola/wag/buffer"
	"github.com/tsavola/wag/disasm"
	"github.com/tsavola/wag/insnmap"
	"github.com/tsavola/wag/internal/test/runner"
	"github.com/tsavola/wag/section"
)

func TestExec(t *testing.T) {
	const (
		filename = "../testdata/exec.wast"

		maxTextSize   = 65536
		maxRODataSize = 4096
		stackSize     = 4096

		dumpBin  = false
		dumpText = false
	)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	wasmReadCloser := wast2wasm(data, false)
	defer wasmReadCloser.Close()
	wasm := bufio.NewReader(wasmReadCloser)

	m := Module{EntrySymbol: "main"}
	m.loadPreliminarySections(wasm, runner.Env)

	var codeBuf bytes.Buffer

	if ok, err := section.CopyCodeSection(&codeBuf, wasm); err != nil {
		t.Fatal(err)
	} else if !ok {
		t.Fatal(ok)
	}

	// skip name section
	if err := section.DiscardUnknownSections(wasm); err != nil {
		t.Fatal(err)
	}

	minMemorySize, maxMemorySize := m.MemoryLimits()

	p, err := runner.NewProgram(maxTextSize, maxRODataSize)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	r, err := p.NewRunner(minMemorySize, maxMemorySize, stackSize)
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	mapping := new(insnmap.Map)

	var printBuf bytes.Buffer
	e, trigger := r.NewExecutor(m.Sigs(), &printBuf)

	m.loadDataSection(wasm, nil)
	p.SetData(m.Data())
	m.loadCodeSection(&codeBuf, buffer.NewFixed(p.Text[:0]), buffer.NewFixed(p.ROData[:0]), p.RODataAddr(), mapping, trigger)
	p.Seal()
	p.SetFuncMap(mapping.FuncAddrs)
	p.SetCallMap(mapping.CallSites)
	if _, err := e.Wait(); err != nil {
		t.Fatal(err)
	}

	if printBuf.Len() > 0 {
		t.Logf("print output:\n%s", string(printBuf.Bytes()))
	}

	if dumpBin {
		if err := writeBin(&m, filename); err != nil {
			t.Error(err)
		}
	}

	if dumpText && testing.Verbose() {
		disasm.Fprint(os.Stdout, m.Text(), mapping.FuncAddrs, nil)
	}
}
