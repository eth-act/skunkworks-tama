// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tamago

package unix

func IsNonblock(fd int) (nonblocking bool, err error) {
	return false, nil
}

func HasNonblockFlag(flag int) bool {
	return false
}
