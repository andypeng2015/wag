// Copyright (c) 2016 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package imports

import (
	"fmt"

	"github.com/tsavola/wag/types"
)

type Function struct {
	types.Function
	Variadic bool
	AbsAddr  uint64
}

func (impl Function) Implements(signature types.Function) bool {
	if impl.Variadic {
		return impl.Function.EqualVariadic(signature)
	} else {
		return impl.Function.Equal(signature)
	}
}

func (f Function) String() (s string) {
	s = fmt.Sprintf("0x%x (", f.AbsAddr)
	for i, t := range f.Args {
		if i > 0 {
			s += ", "
		}
		s += t.String()
	}
	if f.Variadic {
		s += "..."
	}
	s += ") " + f.Result.String()
	return
}
