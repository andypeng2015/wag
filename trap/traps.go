// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package trap

import (
	"fmt"
)

type Id int

const (
	Exit = Id(iota)
	MissingFunction
	Suspended

	CallStackExhausted
	IndirectCallIndex
	IndirectCallSignature
	MemoryOutOfBounds
	Unreachable
	IntegerDivideByZero
	IntegerOverflow

	NumTraps
)

func (id Id) String() string {
	switch id {
	case Exit:
		return "exit"

	case MissingFunction:
		return "missing function"

	case Suspended:
		return "suspended"

	case CallStackExhausted:
		return "call stack exhausted"

	case IndirectCallIndex:
		return "indirect call index out of bounds"

	case IndirectCallSignature:
		return "indirect call signature mismatch"

	case MemoryOutOfBounds:
		return "out of bounds memory access"

	case Unreachable:
		return "unreachable"

	case IntegerDivideByZero:
		return "integer divide by zero"

	case IntegerOverflow:
		return "integer overflow"

	default:
		return fmt.Sprintf("unknown trap %d", id)
	}
}

func (id Id) Error() string {
	return "trap: " + id.String()
}
