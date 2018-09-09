// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package regalloc

import (
	"testing"

	"github.com/tsavola/wag/abi"
	"github.com/tsavola/wag/internal/gen/reg"
	"github.com/tsavola/wag/internal/isa/reglayout"
)

func (a *Allocator) testAllocated(t abi.Type, r reg.R) bool {
	return a.categories[t.Category()].testAllocated(r)
}

func (s *state) testAllocated(r reg.R) bool {
	return s.available&(1<<r) == 0
}

func TestRegAlloc(t *testing.T) {
	a := Make()
	a.CheckNoneAllocated()

	const (
		first = reg.R(reglayout.AllocIntFirst)
		last  = reg.R(reglayout.AllocIntLast)
	)

	for r0 := first; r0 <= last; r0++ {
		r1 := a.AllocResult(abi.I32)
		if r1 == reg.Result {
			t.Fatal(r0, "-", r1, "is result register")
		}
		if r1 != r0 {
			t.Fatal(r1, "is not", r0)
		}

		for r2 := first; r2 <= last; r2++ {
			if a.testAllocated(abi.I32, r2) {
				if r2 > r0 {
					t.Fatal(r0, "-", r2, "allocated")
				}
			} else {
				if r2 <= r0 {
					t.Fatal(r0, "-", r2, "not allocated")
				}
			}
		}
	}

	if r := a.AllocResult(abi.I32); r != reg.Result {
		t.Fatal("allocation succeeded:", r)
	}

	for r0 := first; r0 <= last; r0++ {
		a.Free(abi.I32, r0)
		if a.testAllocated(abi.I32, r0) {
			t.Fatal(r0, "is still allocated")
		}

		r1 := a.AllocResult(abi.I32)
		if r1 == reg.Result {
			t.Fatal("allocation failed")
		}
		if r1 != first {
			t.Fatal(r1, "is not", first)
		}
		a.Free(abi.I32, r1)
	}

	a.CheckNoneAllocated()

	r1 := a.AllocResult(abi.I32)
	if r1 == reg.Result {
		t.Fatal("#1")
	}
	if !a.testAllocated(abi.I32, r1) {
		t.Fatal("#2")
	}
	if r1 != reglayout.AllocIntFirst {
		t.Fatal("#3")
	}

	r2 := a.AllocResult(abi.I32)
	if r2 == reg.Result {
		t.Fatal("#4")
	}
	if !a.testAllocated(abi.I32, r2) {
		t.Fatal("#5")
	}
	if r2 == r1 {
		t.Fatal("#6")
	}
	if r2 != reglayout.AllocIntFirst+1 {
		t.Fatal("#7")
	}

	r1f := a.AllocResult(abi.F32)
	if r1f == reg.Result {
		t.Fatal("#8")
	}
	if !a.testAllocated(abi.F32, r1f) {
		t.Fatal("#9")
	}
	if r1f != reglayout.AllocFloatFirst {
		t.Fatal("#10")
	}
}
