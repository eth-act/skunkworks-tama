// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wasm || windows || tamago

package ld

const syscallExecSupported = false

func (ctxt *Link) execArchive(argv []string) {
	panic("should never arrive here")
}
