// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysrand

import (
	"runtime"
	"sync"
)

var mu sync.Mutex

func read(b []byte) error {
	mu.Lock()
	defer mu.Unlock()

	runtime.GetRandomData(b)

	return nil
}
