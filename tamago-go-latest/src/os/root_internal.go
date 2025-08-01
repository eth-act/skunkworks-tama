// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (js && wasm) || tamago

package os

import (
	"errors"
	"slices"
	"syscall"
)

// checkPathEscapes reports whether name escapes the root.
//
// Due to the lack of openat, checkPathEscapes is subject to TOCTOU races
// when symlinks change during the resolution process.
func checkPathEscapes(r *Root, name string) error {
	return checkPathEscapesInternal(r, name, false)
}

// checkPathEscapesLstat reports whether name escapes the root.
// It does not resolve symlinks in the final path component.
//
// Due to the lack of openat, checkPathEscapes is subject to TOCTOU races
// when symlinks change during the resolution process.
func checkPathEscapesLstat(r *Root, name string) error {
	return checkPathEscapesInternal(r, name, true)
}

func checkPathEscapesInternal(r *Root, name string, lstat bool) error {
	if r.root.closed.Load() {
		return ErrClosed
	}
	parts, suffixSep, err := splitPathInRoot(name, nil, nil)
	if err != nil {
		return err
	}

	i := 0
	symlinks := 0
	base := r.root.name
	for i < len(parts) {
		if parts[i] == ".." {
			// Resolve one or more parent ("..") path components.
			end := i + 1
			for end < len(parts) && parts[end] == ".." {
				end++
			}
			count := end - i
			if count > i {
				return errPathEscapes
			}
			parts = slices.Delete(parts, i-count, end)
			i -= count
			base = r.root.name
			for j := range i {
				base = joinPath(base, parts[j])
			}
			continue
		}

		part := parts[i]
		if i == len(parts)-1 {
			if lstat {
				break
			}
			part += suffixSep
		}

		next := joinPath(base, part)
		fi, err := Lstat(next)
		if err != nil {
			if IsNotExist(err) {
				return nil
			}
			return underlyingError(err)
		}
		if fi.Mode()&ModeSymlink != 0 {
			link, err := Readlink(next)
			if err != nil {
				return errPathEscapes
			}
			symlinks++
			if symlinks > rootMaxSymlinks {
				return errors.New("too many symlinks")
			}
			newparts, newSuffixSep, err := splitPathInRoot(link, parts[:i], parts[i+1:])
			if err != nil {
				return err
			}
			if i == len(parts) {
				// suffixSep contains any trailing path separator characters
				// in the link target.
				// If we are replacing the remainder of the path, retain these.
				// If we're replacing some intermediate component of the path,
				// ignore them, since intermediate components must always be
				// directories.
				suffixSep = newSuffixSep
			}
			parts = newparts
			continue
		}
		if !fi.IsDir() && i < len(parts)-1 {
			return syscall.ENOTDIR
		}

		base = next
		i++
	}
	return nil
}
