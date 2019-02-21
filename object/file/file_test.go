// Copyright (c) 2018 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package file supplies the File type alias, which refers to an OS-dependent
// implementation.  It implements the io.WriterTo interface.  See
// object/file/internal for the common struct specification.
package file

import (
	"io"
	"testing"
)

func TestFileType(*testing.T) {
	var _ io.WriterTo = new(File)
}
