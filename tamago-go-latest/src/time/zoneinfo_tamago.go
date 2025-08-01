// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build tamago

package time

var platformZoneSources = []string{}

func gorootZoneSource(goroot string) (string, bool) {
	return "zoneinfo", true
}

func initLocal() {
	localLoc.name = "UTC"
}
