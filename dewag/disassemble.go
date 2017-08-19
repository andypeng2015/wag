// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dewag

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bnagy/gapstone"

	"github.com/tsavola/wag/sections"
	"github.com/tsavola/wag/traps"
)

func PrintTo(w io.Writer, text, funcMap []byte, ns *sections.NameSection) (err error) {
	var names []sections.FunctionName
	if ns != nil {
		names = ns.FunctionNames
	}

	engine, err := gapstone.New(gapstone.CS_ARCH_X86, gapstone.CS_MODE_64)
	if err != nil {
		return
	}
	defer engine.Close()

	err = engine.SetOption(gapstone.CS_OPT_SYNTAX, gapstone.CS_OPT_SYNTAX_ATT)
	if err != nil {
		return
	}

	insns, err := engine.Disasm(text, 0, 0)
	if err != nil {
		return
	}

	targets := map[uint]string{
		16: "start",
	}

	firstFuncAddr := uint(binary.LittleEndian.Uint32(funcMap))

	for i := 0; len(funcMap) > 0; i++ {
		addr := binary.LittleEndian.Uint32(funcMap)
		funcMap = funcMap[4:]

		var name string
		if i < len(names) {
			name = names[i].FunName
		} else {
			name = fmt.Sprintf("func_%d", i)
		}

		targets[uint(addr)] = name
	}

	sequence := 0

	for i := range insns {
		insn := insns[i]
		var name string

		switch {
		case insn.Mnemonic == "jmpq":
			continue

		case strings.HasPrefix(insn.Mnemonic, "j"):

		case insn.Mnemonic == "callq" && strings.HasPrefix(insn.OpStr, "0x"):

		case insn.Address < firstFuncAddr && insn.Mnemonic == "movl" && strings.HasPrefix(insn.OpStr, "$") && strings.HasSuffix(insn.OpStr, ", %eax"):
			var n uint
			fmt.Sscanf(insn.OpStr, "$%d, %%eax", &n)
			if id := traps.Id(n); id < traps.NumTraps {
				targets[insn.Address] = strings.Replace(id.String(), " ", "_", -1)
			}
			continue

		default:
			continue
		}

		addr, err := strconv.ParseUint(insn.OpStr, 0, 32)
		if err != nil {
			panic(err)
		}

		name, found := targets[uint(addr)]
		if !found {
			name = fmt.Sprintf(".L%d", sequence)
			sequence++

			targets[uint(addr)] = name
		}

		insns[i].OpStr = name
	}

	skip := false

	for _, insn := range insns {
		name, found := targets[insn.Address]
		if found {
			if !strings.HasPrefix(name, ".") && name != "exit" {
				fmt.Fprintln(w)
			}
			fmt.Fprintf(w, "%s:\n", name)
		}

		switch insn.Id {
		case gapstone.X86_INS_INT3:
			if skip {
				continue
			}
			skip = true
			fallthrough

		default:
			fmt.Fprintf(w, "\t%s\t%s\n", insn.Mnemonic, insn.OpStr)
		}
	}

	fmt.Fprintln(w)
	return
}
