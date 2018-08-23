// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package meta

type TextAddr int32

type CallSite struct {
	ReturnAddr  TextAddr
	StackOffset int32
}
