// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build wagarm64 arm64,!wagamd64

package reglayout

const (
	AllocIntFirst = 2
	AllocIntLast  = 25

	AllocFloatFirst = 2
	AllocFloatLast  = 30 // TODO: is there a zero reg at 31?

	Radix = 32
)
