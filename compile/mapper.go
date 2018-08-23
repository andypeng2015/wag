// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compile

import (
	"github.com/tsavola/wag/internal/module"
	"github.com/tsavola/wag/meta"
)

type Mapper = module.Mapper

type dummyMapper struct{}

func (dummyMapper) InitModule(int, int)          {}
func (dummyMapper) PutImportFunc(meta.TextAddr)  {}
func (dummyMapper) PutFunc(meta.TextAddr)        {}
func (dummyMapper) PutCall(meta.TextAddr, int32) {}
func (dummyMapper) PutInsn(meta.TextAddr)        {}
