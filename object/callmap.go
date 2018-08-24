// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package object

// CallSite represents a position within the text section (machine code) where
// a function call is made.
type CallSite struct {
	ReturnAddr  int32 // The address immediately after the call instruction
	StackOffset int32 // Calling function's stack usage at time of call
}

// CallMap implements compile.ObjectMap.  It stores all function addresses and
// call sites, but no instruction information.
type CallMap struct {
	FuncMap
	CallSites []CallSite
}

func (m *CallMap) InitObjectMap(numImportFuncs, numOtherFuncs int) {
	m.FuncMap.InitObjectMap(numImportFuncs, numOtherFuncs)
	m.CallSites = make([]CallSite, 0, numImportFuncs+numOtherFuncs) // conservative guess
}

func (m *CallMap) PutCallSite(retAddr int32, stackOffset int32) {
	m.CallSites = append(m.CallSites, CallSite{retAddr, stackOffset})
}
