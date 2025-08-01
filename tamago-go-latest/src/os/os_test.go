// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"internal/testenv"
	"io"
	"io/fs"
	"log"
	. "os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strings"
	"sync"
	"syscall"
	"testing"
	"testing/fstest"
	"time"
)

func TestMain(m *testing.M) {
	if Getenv("GO_OS_TEST_DRAIN_STDIN") == "1" {
		Stdout.Close()
		io.Copy(io.Discard, Stdin)
		Exit(0)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	Exit(m.Run())
}

var dot = []string{
	"dir_unix.go",
	"env.go",
	"error.go",
	"file.go",
	"os_test.go",
	"types.go",
	"stat_darwin.go",
	"stat_linux.go",
}

type sysDir struct {
	name  string
	files []string
}

var sysdir = func() *sysDir {
	switch runtime.GOOS {
	case "android":
		return &sysDir{
			"/system/lib",
			[]string{
				"libmedia.so",
				"libpowermanager.so",
			},
		}
	case "ios":
		wd, err := syscall.Getwd()
		if err != nil {
			wd = err.Error()
		}
		sd := &sysDir{
			filepath.Join(wd, "..", ".."),
			[]string{
				"ResourceRules.plist",
				"Info.plist",
			},
		}
		found := true
		for _, f := range sd.files {
			path := filepath.Join(sd.name, f)
			if _, err := Stat(path); err != nil {
				found = false
				break
			}
		}
		if found {
			return sd
		}
		// In a self-hosted iOS build the above files might
		// not exist. Look for system files instead below.
	case "windows":
		return &sysDir{
			Getenv("SystemRoot") + "\\system32\\drivers\\etc",
			[]string{
				"networks",
				"protocol",
				"services",
			},
		}
	case "plan9":
		return &sysDir{
			"/lib/ndb",
			[]string{
				"common",
				"local",
			},
		}
	case "wasip1":
		// wasmtime has issues resolving symbolic links that are often present
		// in directories like /etc/group below (e.g. private/etc/group on OSX).
		// For this reason we use files in the Go source tree instead.
		return &sysDir{
			runtime.GOROOT(),
			[]string{
				"go.env",
				"LICENSE",
				"CONTRIBUTING.md",
			},
		}
	case "tamago":
		return &sysDir{
			"testdata",
			[]string{
				"hello",
			},
		}
	}
	return &sysDir{
		"/etc",
		[]string{
			"group",
			"hosts",
			"passwd",
		},
	}
}()

func size(name string, t *testing.T) int64 {
	file, err := Open(name)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Error(err)
		}
	}()
	n, err := io.Copy(io.Discard, file)
	if err != nil {
		t.Fatal(err)
	}
	return n
}

func equal(name1, name2 string) (r bool) {
	switch runtime.GOOS {
	case "windows":
		r = strings.EqualFold(name1, name2)
	default:
		r = name1 == name2
	}
	return
}

func newFile(t *testing.T) (f *File) {
	t.Helper()
	f, err := CreateTemp("", "_Go_"+t.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := f.Close(); err != nil && !errors.Is(err, ErrClosed) {
			t.Fatal(err)
		}
		if err := Remove(f.Name()); err != nil {
			t.Fatal(err)
		}
	})
	return
}

var sfdir = sysdir.name
var sfname = sysdir.files[0]

func TestStat(t *testing.T) {
	t.Parallel()

	path := sfdir + "/" + sfname
	dir, err := Stat(path)
	if err != nil {
		t.Fatal("stat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

func TestStatError(t *testing.T) {
	t.Chdir(t.TempDir())

	path := "no-such-file"

	fi, err := Stat(path)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if fi != nil {
		t.Errorf("got %v, want nil", fi)
	}
	if perr, ok := err.(*PathError); !ok {
		t.Errorf("got %T, want %T", err, perr)
	}

	testenv.MustHaveSymlink(t)

	link := "symlink"
	err = Symlink(path, link)
	if err != nil {
		t.Fatal(err)
	}

	fi, err = Stat(link)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if fi != nil {
		t.Errorf("got %v, want nil", fi)
	}
	if perr, ok := err.(*PathError); !ok {
		t.Errorf("got %T, want %T", err, perr)
	}
}

func TestStatSymlinkLoop(t *testing.T) {
	testenv.MustHaveSymlink(t)
	t.Chdir(t.TempDir())

	err := Symlink("x", "y")
	if err != nil {
		t.Fatal(err)
	}
	defer Remove("y")

	err = Symlink("y", "x")
	if err != nil {
		t.Fatal(err)
	}
	defer Remove("x")

	_, err = Stat("x")
	if _, ok := err.(*fs.PathError); !ok {
		t.Errorf("expected *PathError, got %T: %v\n", err, err)
	}
}

func TestFstat(t *testing.T) {
	t.Parallel()

	path := sfdir + "/" + sfname
	file, err1 := Open(path)
	if err1 != nil {
		t.Fatal("open failed:", err1)
	}
	defer file.Close()
	dir, err2 := file.Stat()
	if err2 != nil {
		t.Fatal("fstat failed:", err2)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	filesize := size(path, t)
	if dir.Size() != filesize {
		t.Error("size should be", filesize, "; is", dir.Size())
	}
}

func TestLstat(t *testing.T) {
	t.Parallel()

	path := sfdir + "/" + sfname
	dir, err := Lstat(path)
	if err != nil {
		t.Fatal("lstat failed:", err)
	}
	if !equal(sfname, dir.Name()) {
		t.Error("name should be ", sfname, "; is", dir.Name())
	}
	if dir.Mode()&ModeSymlink == 0 {
		filesize := size(path, t)
		if dir.Size() != filesize {
			t.Error("size should be", filesize, "; is", dir.Size())
		}
	}
}

// Read with length 0 should not return EOF.
func TestRead0(t *testing.T) {
	t.Parallel()

	path := sfdir + "/" + sfname
	f, err := Open(path)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	defer f.Close()

	b := make([]byte, 0)
	n, err := f.Read(b)
	if n != 0 || err != nil {
		t.Errorf("Read(0) = %d, %v, want 0, nil", n, err)
	}
	b = make([]byte, 100)
	n, err = f.Read(b)
	if n <= 0 || err != nil {
		t.Errorf("Read(100) = %d, %v, want >0, nil", n, err)
	}
}

// Reading a closed file should return ErrClosed error
func TestReadClosed(t *testing.T) {
	t.Parallel()

	path := sfdir + "/" + sfname
	file, err := Open(path)
	if err != nil {
		t.Fatal("open failed:", err)
	}
	file.Close() // close immediately

	b := make([]byte, 100)
	_, err = file.Read(b)

	e, ok := err.(*PathError)
	if !ok || e.Err != ErrClosed {
		t.Fatalf("Read: got %T(%v), want %T(%v)", err, err, e, ErrClosed)
	}
}

func testReaddirnames(dir string, contents []string) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		file, err := Open(dir)
		if err != nil {
			t.Fatalf("open %q failed: %v", dir, err)
		}
		defer file.Close()
		s, err2 := file.Readdirnames(-1)
		if err2 != nil {
			t.Fatalf("Readdirnames %q failed: %v", dir, err2)
		}
		for _, m := range contents {
			found := false
			for _, n := range s {
				if n == "." || n == ".." {
					t.Errorf("got %q in directory", n)
				}
				if !equal(m, n) {
					continue
				}
				if found {
					t.Error("present twice:", m)
				}
				found = true
			}
			if !found {
				t.Error("could not find", m)
			}
		}
		if s == nil {
			t.Error("Readdirnames returned nil instead of empty slice")
		}
	}
}

func testReaddir(dir string, contents []string) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		file, err := Open(dir)
		if err != nil {
			t.Fatalf("open %q failed: %v", dir, err)
		}
		defer file.Close()
		s, err2 := file.Readdir(-1)
		if err2 != nil {
			t.Fatalf("Readdir %q failed: %v", dir, err2)
		}
		for _, m := range contents {
			found := false
			for _, n := range s {
				if n.Name() == "." || n.Name() == ".." {
					t.Errorf("got %q in directory", n.Name())
				}
				if !equal(m, n.Name()) {
					continue
				}
				if found {
					t.Error("present twice:", m)
				}
				found = true
			}
			if !found {
				t.Error("could not find", m)
			}
		}
		if s == nil {
			t.Error("Readdir returned nil instead of empty slice")
		}
	}
}

func testReadDir(dir string, contents []string) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		file, err := Open(dir)
		if err != nil {
			t.Fatalf("open %q failed: %v", dir, err)
		}
		defer file.Close()
		s, err2 := file.ReadDir(-1)
		if err2 != nil {
			t.Fatalf("ReadDir %q failed: %v", dir, err2)
		}
		for _, m := range contents {
			found := false
			for _, n := range s {
				if n.Name() == "." || n.Name() == ".." {
					t.Errorf("got %q in directory", n)
				}
				if !equal(m, n.Name()) {
					continue
				}
				if found {
					t.Error("present twice:", m)
				}
				found = true
				lstat, err := Lstat(dir + "/" + m)
				if err != nil {
					t.Fatal(err)
				}
				if n.IsDir() != lstat.IsDir() {
					t.Errorf("%s: IsDir=%v, want %v", m, n.IsDir(), lstat.IsDir())
				}
				if n.Type() != lstat.Mode().Type() {
					t.Errorf("%s: IsDir=%v, want %v", m, n.Type(), lstat.Mode().Type())
				}
				info, err := n.Info()
				if err != nil {
					t.Errorf("%s: Info: %v", m, err)
					continue
				}
				if !SameFile(info, lstat) {
					t.Errorf("%s: Info: SameFile(info, lstat) = false", m)
				}
			}
			if !found {
				t.Error("could not find", m)
			}
		}
		if s == nil {
			t.Error("ReadDir returned nil instead of empty slice")
		}
	}
}

func TestFileReaddirnames(t *testing.T) {
	t.Parallel()

	t.Run(".", testReaddirnames(".", dot))
	t.Run("sysdir", testReaddirnames(sysdir.name, sysdir.files))
	t.Run("TempDir", testReaddirnames(t.TempDir(), nil))
}

func TestFileReaddir(t *testing.T) {
	t.Parallel()

	t.Run(".", testReaddir(".", dot))
	t.Run("sysdir", testReaddir(sysdir.name, sysdir.files))
	t.Run("TempDir", testReaddir(t.TempDir(), nil))
}

func TestFileReadDir(t *testing.T) {
	t.Parallel()

	t.Run(".", testReadDir(".", dot))
	t.Run("sysdir", testReadDir(sysdir.name, sysdir.files))
	t.Run("TempDir", testReadDir(t.TempDir(), nil))
}

func benchmarkReaddirname(path string, b *testing.B) {
	var nentries int
	for i := 0; i < b.N; i++ {
		f, err := Open(path)
		if err != nil {
			b.Fatalf("open %q failed: %v", path, err)
		}
		ns, err := f.Readdirnames(-1)
		f.Close()
		if err != nil {
			b.Fatalf("readdirnames %q failed: %v", path, err)
		}
		nentries = len(ns)
	}
	b.Logf("benchmarkReaddirname %q: %d entries", path, nentries)
}

func benchmarkReaddir(path string, b *testing.B) {
	var nentries int
	for i := 0; i < b.N; i++ {
		f, err := Open(path)
		if err != nil {
			b.Fatalf("open %q failed: %v", path, err)
		}
		fs, err := f.Readdir(-1)
		f.Close()
		if err != nil {
			b.Fatalf("readdir %q failed: %v", path, err)
		}
		nentries = len(fs)
	}
	b.Logf("benchmarkReaddir %q: %d entries", path, nentries)
}

func benchmarkReadDir(path string, b *testing.B) {
	var nentries int
	for i := 0; i < b.N; i++ {
		f, err := Open(path)
		if err != nil {
			b.Fatalf("open %q failed: %v", path, err)
		}
		fs, err := f.ReadDir(-1)
		f.Close()
		if err != nil {
			b.Fatalf("readdir %q failed: %v", path, err)
		}
		nentries = len(fs)
	}
	b.Logf("benchmarkReadDir %q: %d entries", path, nentries)
}

func BenchmarkReaddirname(b *testing.B) {
	benchmarkReaddirname(".", b)
}

func BenchmarkReaddir(b *testing.B) {
	benchmarkReaddir(".", b)
}

func BenchmarkReadDir(b *testing.B) {
	benchmarkReadDir(".", b)
}

func benchmarkStat(b *testing.B, path string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Stat(path)
		if err != nil {
			b.Fatalf("Stat(%q) failed: %v", path, err)
		}
	}
}

func benchmarkLstat(b *testing.B, path string) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Lstat(path)
		if err != nil {
			b.Fatalf("Lstat(%q) failed: %v", path, err)
		}
	}
}

func BenchmarkStatDot(b *testing.B) {
	benchmarkStat(b, ".")
}

func BenchmarkStatFile(b *testing.B) {
	benchmarkStat(b, filepath.Join(runtime.GOROOT(), "src/os/os_test.go"))
}

func BenchmarkStatDir(b *testing.B) {
	benchmarkStat(b, filepath.Join(runtime.GOROOT(), "src/os"))
}

func BenchmarkLstatDot(b *testing.B) {
	benchmarkLstat(b, ".")
}

func BenchmarkLstatFile(b *testing.B) {
	benchmarkLstat(b, filepath.Join(runtime.GOROOT(), "src/os/os_test.go"))
}

func BenchmarkLstatDir(b *testing.B) {
	benchmarkLstat(b, filepath.Join(runtime.GOROOT(), "src/os"))
}

// Read the directory one entry at a time.
func smallReaddirnames(file *File, length int, t *testing.T) []string {
	names := make([]string, length)
	count := 0
	for {
		d, err := file.Readdirnames(1)
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("readdirnames %q failed: %v", file.Name(), err)
		}
		if len(d) == 0 {
			t.Fatalf("readdirnames %q returned empty slice and no error", file.Name())
		}
		names[count] = d[0]
		count++
	}
	return names[0:count]
}

// Check that reading a directory one entry at a time gives the same result
// as reading it all at once.
func TestReaddirnamesOneAtATime(t *testing.T) {
	t.Parallel()

	// big directory that doesn't change often.
	dir := "/usr/bin"
	switch runtime.GOOS {
	case "android":
		dir = "/system/bin"
	case "ios", "wasip1", "tamago":
		wd, err := Getwd()
		if err != nil {
			t.Fatal(err)
		}
		dir = wd
	case "plan9":
		dir = "/bin"
	case "windows":
		dir = Getenv("SystemRoot") + "\\system32"
	}
	file, err := Open(dir)
	if err != nil {
		t.Fatalf("open %q failed: %v", dir, err)
	}
	defer file.Close()
	all, err1 := file.Readdirnames(-1)
	if err1 != nil {
		t.Fatalf("readdirnames %q failed: %v", dir, err1)
	}
	file1, err2 := Open(dir)
	if err2 != nil {
		t.Fatalf("open %q failed: %v", dir, err2)
	}
	defer file1.Close()
	small := smallReaddirnames(file1, len(all)+100, t) // +100 in case we screw up
	if len(small) < len(all) {
		t.Fatalf("len(small) is %d, less than %d", len(small), len(all))
	}
	for i, n := range all {
		if small[i] != n {
			t.Errorf("small read %q mismatch: %v", small[i], n)
		}
	}
}

func TestReaddirNValues(t *testing.T) {
	if testing.Short() {
		t.Skip("test.short; skipping")
	}
	t.Parallel()

	dir := t.TempDir()
	for i := 1; i <= 105; i++ {
		f, err := Create(filepath.Join(dir, fmt.Sprintf("%d", i)))
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		f.Write([]byte(strings.Repeat("X", i)))
		f.Close()
	}

	var d *File
	openDir := func() {
		var err error
		d, err = Open(dir)
		if err != nil {
			t.Fatalf("Open directory: %v", err)
		}
	}

	readdirExpect := func(n, want int, wantErr error) {
		t.Helper()
		fi, err := d.Readdir(n)
		if err != wantErr {
			t.Fatalf("Readdir of %d got error %v, want %v", n, err, wantErr)
		}
		if g, e := len(fi), want; g != e {
			t.Errorf("Readdir of %d got %d files, want %d", n, g, e)
		}
	}

	readDirExpect := func(n, want int, wantErr error) {
		t.Helper()
		de, err := d.ReadDir(n)
		if err != wantErr {
			t.Fatalf("ReadDir of %d got error %v, want %v", n, err, wantErr)
		}
		if g, e := len(de), want; g != e {
			t.Errorf("ReadDir of %d got %d files, want %d", n, g, e)
		}
	}

	readdirnamesExpect := func(n, want int, wantErr error) {
		t.Helper()
		fi, err := d.Readdirnames(n)
		if err != wantErr {
			t.Fatalf("Readdirnames of %d got error %v, want %v", n, err, wantErr)
		}
		if g, e := len(fi), want; g != e {
			t.Errorf("Readdirnames of %d got %d files, want %d", n, g, e)
		}
	}

	for _, fn := range []func(int, int, error){readdirExpect, readdirnamesExpect, readDirExpect} {
		// Test the slurp case
		openDir()
		fn(0, 105, nil)
		fn(0, 0, nil)
		d.Close()

		// Slurp with -1 instead
		openDir()
		fn(-1, 105, nil)
		fn(-2, 0, nil)
		fn(0, 0, nil)
		d.Close()

		// Test the bounded case
		openDir()
		fn(1, 1, nil)
		fn(2, 2, nil)
		fn(105, 102, nil) // and tests buffer >100 case
		fn(3, 0, io.EOF)
		d.Close()
	}
}

func touch(t *testing.T, name string) {
	f, err := Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestReaddirStatFailures(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "plan9":
		// Windows and Plan 9 already do this correctly,
		// but are structured with different syscalls such
		// that they don't use Lstat, so the hook below for
		// testing it wouldn't work.
		t.Skipf("skipping test on %v", runtime.GOOS)
	}

	var xerr error // error to return for x
	*LstatP = func(path string) (FileInfo, error) {
		if xerr != nil && strings.HasSuffix(path, "x") {
			return nil, xerr
		}
		return Lstat(path)
	}
	defer func() { *LstatP = Lstat }()

	dir := t.TempDir()
	touch(t, filepath.Join(dir, "good1"))
	touch(t, filepath.Join(dir, "x")) // will disappear or have an error
	touch(t, filepath.Join(dir, "good2"))
	readDir := func() ([]FileInfo, error) {
		d, err := Open(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer d.Close()
		return d.Readdir(-1)
	}
	mustReadDir := func(testName string) []FileInfo {
		fis, err := readDir()
		if err != nil {
			t.Fatalf("%s: Readdir: %v", testName, err)
		}
		return fis
	}
	names := func(fis []FileInfo) []string {
		s := make([]string, len(fis))
		for i, fi := range fis {
			s[i] = fi.Name()
		}
		slices.Sort(s)
		return s
	}

	if got, want := names(mustReadDir("initial readdir")),
		[]string{"good1", "good2", "x"}; !slices.Equal(got, want) {
		t.Errorf("initial readdir got %q; want %q", got, want)
	}

	xerr = ErrNotExist
	if got, want := names(mustReadDir("with x disappearing")),
		[]string{"good1", "good2"}; !slices.Equal(got, want) {
		t.Errorf("with x disappearing, got %q; want %q", got, want)
	}

	xerr = errors.New("some real error")
	if _, err := readDir(); err != xerr {
		t.Errorf("with a non-ErrNotExist error, got error %v; want %v", err, xerr)
	}
}

// Readdir on a regular file should fail.
func TestReaddirOfFile(t *testing.T) {
	t.Parallel()

	f, err := CreateTemp(t.TempDir(), "_Go_ReaddirOfFile")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte("foo"))
	f.Close()
	reg, err := Open(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer reg.Close()

	names, err := reg.Readdirnames(-1)
	if err == nil {
		t.Error("Readdirnames succeeded; want non-nil error")
	}
	var pe *PathError
	if !errors.As(err, &pe) || pe.Path != f.Name() {
		t.Errorf("Readdirnames returned %q; want a PathError with path %q", err, f.Name())
	}
	if len(names) > 0 {
		t.Errorf("unexpected dir names in regular file: %q", names)
	}
}

func TestHardLink(t *testing.T) {
	testenv.MustHaveLink(t)
	t.Chdir(t.TempDir())

	from, to := "hardlinktestfrom", "hardlinktestto"
	file, err := Create(to)
	if err != nil {
		t.Fatalf("open %q failed: %v", to, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("close %q failed: %v", to, err)
	}
	err = Link(to, from)
	if err != nil {
		t.Fatalf("link %q, %q failed: %v", to, from, err)
	}

	none := "hardlinktestnone"
	err = Link(none, none)
	// Check the returned error is well-formed.
	if lerr, ok := err.(*LinkError); !ok || lerr.Error() == "" {
		t.Errorf("link %q, %q failed to return a valid error", none, none)
	}

	tostat, err := Stat(to)
	if err != nil {
		t.Fatalf("stat %q failed: %v", to, err)
	}
	fromstat, err := Stat(from)
	if err != nil {
		t.Fatalf("stat %q failed: %v", from, err)
	}
	if !SameFile(tostat, fromstat) {
		t.Errorf("link %q, %q did not create hard link", to, from)
	}
	// We should not be able to perform the same Link() a second time
	err = Link(to, from)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "link" {
			t.Errorf("Link(%q, %q) err.Op = %q; want %q", to, from, err.Op, "link")
		}
		if err.Old != to {
			t.Errorf("Link(%q, %q) err.Old = %q; want %q", to, from, err.Old, to)
		}
		if err.New != from {
			t.Errorf("Link(%q, %q) err.New = %q; want %q", to, from, err.New, from)
		}
		if !IsExist(err.Err) {
			t.Errorf("Link(%q, %q) err.Err = %q; want %q", to, from, err.Err, "file exists error")
		}
	case nil:
		t.Errorf("link %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("link %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

func TestSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)
	t.Chdir(t.TempDir())

	from, to := "symlinktestfrom", "symlinktestto"
	file, err := Create(to)
	if err != nil {
		t.Fatalf("Create(%q) failed: %v", to, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("Close(%q) failed: %v", to, err)
	}
	err = Symlink(to, from)
	if err != nil {
		t.Fatalf("Symlink(%q, %q) failed: %v", to, from, err)
	}
	tostat, err := Lstat(to)
	if err != nil {
		t.Fatalf("Lstat(%q) failed: %v", to, err)
	}
	if tostat.Mode()&ModeSymlink != 0 {
		t.Fatalf("Lstat(%q).Mode()&ModeSymlink = %v, want 0", to, tostat.Mode()&ModeSymlink)
	}
	fromstat, err := Stat(from)
	if err != nil {
		t.Fatalf("Stat(%q) failed: %v", from, err)
	}
	if !SameFile(tostat, fromstat) {
		t.Errorf("Symlink(%q, %q) did not create symlink", to, from)
	}
	fromstat, err = Lstat(from)
	if err != nil {
		t.Fatalf("Lstat(%q) failed: %v", from, err)
	}
	if fromstat.Mode()&ModeSymlink == 0 {
		t.Fatalf("Lstat(%q).Mode()&ModeSymlink = 0, want %v", from, ModeSymlink)
	}
	fromstat, err = Stat(from)
	if err != nil {
		t.Fatalf("Stat(%q) failed: %v", from, err)
	}
	if fromstat.Name() != from {
		t.Errorf("Stat(%q).Name() = %q, want %q", from, fromstat.Name(), from)
	}
	if fromstat.Mode()&ModeSymlink != 0 {
		t.Fatalf("Stat(%q).Mode()&ModeSymlink = %v, want 0", from, fromstat.Mode()&ModeSymlink)
	}
	s, err := Readlink(from)
	if err != nil {
		t.Fatalf("Readlink(%q) failed: %v", from, err)
	}
	if s != to {
		t.Fatalf("Readlink(%q) = %q, want %q", from, s, to)
	}
	file, err = Open(from)
	if err != nil {
		t.Fatalf("Open(%q) failed: %v", from, err)
	}
	file.Close()
}

func TestLongSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)
	t.Chdir(t.TempDir())

	s := "0123456789abcdef"
	// Long, but not too long: a common limit is 255.
	s = s + s + s + s + s + s + s + s + s + s + s + s + s + s + s
	from := "longsymlinktestfrom"
	err := Symlink(s, from)
	if err != nil {
		t.Fatalf("symlink %q, %q failed: %v", s, from, err)
	}
	r, err := Readlink(from)
	if err != nil {
		t.Fatalf("readlink %q failed: %v", from, err)
	}
	if r != s {
		t.Fatalf("after symlink %q != %q", r, s)
	}
}

func TestRename(t *testing.T) {
	t.Chdir(t.TempDir())
	from, to := "renamefrom", "renameto"

	file, err := Create(from)
	if err != nil {
		t.Fatalf("open %q failed: %v", from, err)
	}
	if err = file.Close(); err != nil {
		t.Errorf("close %q failed: %v", from, err)
	}
	err = Rename(from, to)
	if err != nil {
		t.Fatalf("rename %q, %q failed: %v", to, from, err)
	}
	_, err = Stat(to)
	if err != nil {
		t.Errorf("stat %q failed: %v", to, err)
	}
}

func TestRenameOverwriteDest(t *testing.T) {
	t.Chdir(t.TempDir())
	from, to := "renamefrom", "renameto"

	toData := []byte("to")
	fromData := []byte("from")

	err := WriteFile(to, toData, 0777)
	if err != nil {
		t.Fatalf("write file %q failed: %v", to, err)
	}

	err = WriteFile(from, fromData, 0777)
	if err != nil {
		t.Fatalf("write file %q failed: %v", from, err)
	}
	err = Rename(from, to)
	if err != nil {
		t.Fatalf("rename %q, %q failed: %v", to, from, err)
	}

	_, err = Stat(from)
	if err == nil {
		t.Errorf("from file %q still exists", from)
	}
	if err != nil && !IsNotExist(err) {
		t.Fatalf("stat from: %v", err)
	}
	toFi, err := Stat(to)
	if err != nil {
		t.Fatalf("stat %q failed: %v", to, err)
	}
	if toFi.Size() != int64(len(fromData)) {
		t.Errorf(`"to" size = %d; want %d (old "from" size)`, toFi.Size(), len(fromData))
	}
}

func TestRenameFailed(t *testing.T) {
	t.Chdir(t.TempDir())
	from, to := "renamefrom", "renameto"

	err := Rename(from, to)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "rename" {
			t.Errorf("rename %q, %q: err.Op: want %q, got %q", from, to, "rename", err.Op)
		}
		if err.Old != from {
			t.Errorf("rename %q, %q: err.Old: want %q, got %q", from, to, from, err.Old)
		}
		if err.New != to {
			t.Errorf("rename %q, %q: err.New: want %q, got %q", from, to, to, err.New)
		}
	case nil:
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("rename %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

func TestRenameNotExisting(t *testing.T) {
	t.Chdir(t.TempDir())
	from, to := "doesnt-exist", "dest"

	Mkdir(to, 0777)

	if err := Rename(from, to); !IsNotExist(err) {
		t.Errorf("Rename(%q, %q) = %v; want an IsNotExist error", from, to, err)
	}
}

func TestRenameToDirFailed(t *testing.T) {
	t.Chdir(t.TempDir())
	from, to := "renamefrom", "renameto"

	Mkdir(from, 0777)
	Mkdir(to, 0777)

	err := Rename(from, to)
	switch err := err.(type) {
	case *LinkError:
		if err.Op != "rename" {
			t.Errorf("rename %q, %q: err.Op: want %q, got %q", from, to, "rename", err.Op)
		}
		if err.Old != from {
			t.Errorf("rename %q, %q: err.Old: want %q, got %q", from, to, from, err.Old)
		}
		if err.New != to {
			t.Errorf("rename %q, %q: err.New: want %q, got %q", from, to, to, err.New)
		}
	case nil:
		t.Errorf("rename %q, %q: expected error, got nil", from, to)
	default:
		t.Errorf("rename %q, %q: expected %T, got %T %v", from, to, new(LinkError), err, err)
	}
}

func TestRenameCaseDifference(pt *testing.T) {
	from, to := "renameFROM", "RENAMEfrom"
	tests := []struct {
		name   string
		create func() error
	}{
		{"dir", func() error {
			return Mkdir(from, 0777)
		}},
		{"file", func() error {
			fd, err := Create(from)
			if err != nil {
				return err
			}
			return fd.Close()
		}},
	}

	for _, test := range tests {
		pt.Run(test.name, func(t *testing.T) {
			t.Chdir(t.TempDir())

			if err := test.create(); err != nil {
				t.Fatalf("failed to create test file: %s", err)
			}

			if _, err := Stat(to); err != nil {
				// Sanity check that the underlying filesystem is not case sensitive.
				if IsNotExist(err) {
					t.Skipf("case sensitive filesystem")
				}
				t.Fatalf("stat %q, got: %q", to, err)
			}

			if err := Rename(from, to); err != nil {
				t.Fatalf("unexpected error when renaming from %q to %q: %s", from, to, err)
			}

			fd, err := Open(".")
			if err != nil {
				t.Fatalf("Open .: %s", err)
			}

			// Stat does not return the real case of the file (it returns what the called asked for)
			// So we have to use readdir to get the real name of the file.
			dirNames, err := fd.Readdirnames(-1)
			fd.Close()
			if err != nil {
				t.Fatalf("readdirnames: %s", err)
			}

			if dirNamesLen := len(dirNames); dirNamesLen != 1 {
				t.Fatalf("unexpected dirNames len, got %q, want %q", dirNamesLen, 1)
			}

			if dirNames[0] != to {
				t.Errorf("unexpected name, got %q, want %q", dirNames[0], to)
			}
		})
	}
}

func testStartProcess(dir, cmd string, args []string, expect string) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		r, w, err := Pipe()
		if err != nil {
			t.Fatalf("Pipe: %v", err)
		}
		defer r.Close()
		attr := &ProcAttr{Dir: dir, Files: []*File{nil, w, Stderr}}
		p, err := StartProcess(cmd, args, attr)
		if err != nil {
			t.Fatalf("StartProcess: %v", err)
		}
		w.Close()

		var b strings.Builder
		io.Copy(&b, r)
		output := b.String()

		fi1, _ := Stat(strings.TrimSpace(output))
		fi2, _ := Stat(expect)
		if !SameFile(fi1, fi2) {
			t.Errorf("exec %q returned %q wanted %q",
				strings.Join(append([]string{cmd}, args...), " "), output, expect)
		}
		p.Wait()
	}
}

func TestStartProcess(t *testing.T) {
	testenv.MustHaveExec(t)
	t.Parallel()

	var dir, cmd string
	var args []string
	switch runtime.GOOS {
	case "android":
		t.Skip("android doesn't have /bin/pwd")
	case "windows":
		cmd = Getenv("COMSPEC")
		dir = Getenv("SystemRoot")
		args = []string{"/c", "cd"}
	default:
		var err error
		cmd, err = exec.LookPath("pwd")
		if err != nil {
			t.Fatalf("Can't find pwd: %v", err)
		}
		dir = "/"
		args = []string{}
		t.Logf("Testing with %v", cmd)
	}
	cmddir, cmdbase := filepath.Split(cmd)
	args = append([]string{cmdbase}, args...)
	t.Run("absolute", testStartProcess(dir, cmd, args, dir))
	t.Run("relative", testStartProcess(cmddir, cmdbase, args, cmddir))
}

func checkMode(t *testing.T, path string, mode FileMode) {
	dir, err := Stat(path)
	if err != nil {
		t.Fatalf("Stat %q (looking for mode %#o): %s", path, mode, err)
	}
	if dir.Mode()&ModePerm != mode {
		t.Errorf("Stat %q: mode %#o want %#o", path, dir.Mode(), mode)
	}
}

func TestChmod(t *testing.T) {
	// Chmod is not supported on wasip1.
	if runtime.GOOS == "wasip1" {
		t.Skip("Chmod is not supported on " + runtime.GOOS)
	}
	t.Parallel()

	f := newFile(t)
	// Creation mode is read write

	fm := FileMode(0456)
	if runtime.GOOS == "windows" {
		fm = FileMode(0444) // read-only file
	}
	if err := Chmod(f.Name(), fm); err != nil {
		t.Fatalf("chmod %s %#o: %s", f.Name(), fm, err)
	}
	checkMode(t, f.Name(), fm)

	fm = FileMode(0123)
	if runtime.GOOS == "windows" {
		fm = FileMode(0666) // read-write file
	}
	if err := f.Chmod(fm); err != nil {
		t.Fatalf("chmod %s %#o: %s", f.Name(), fm, err)
	}
	checkMode(t, f.Name(), fm)
}

func checkSize(t *testing.T, f *File, size int64) {
	t.Helper()
	dir, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat %q (looking for size %d): %s", f.Name(), size, err)
	}
	if dir.Size() != size {
		t.Errorf("Stat %q: size %d want %d", f.Name(), dir.Size(), size)
	}
}

func TestFTruncate(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	checkSize(t, f, 0)
	f.Write([]byte("hello, world\n"))
	checkSize(t, f, 13)
	f.Truncate(10)
	checkSize(t, f, 10)
	f.Truncate(1024)
	checkSize(t, f, 1024)
	f.Truncate(0)
	checkSize(t, f, 0)
	_, err := f.Write([]byte("surprise!"))
	if err == nil {
		checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	checkSize(t, f, 0)
	f.Write([]byte("hello, world\n"))
	checkSize(t, f, 13)
	Truncate(f.Name(), 10)
	checkSize(t, f, 10)
	Truncate(f.Name(), 1024)
	checkSize(t, f, 1024)
	Truncate(f.Name(), 0)
	checkSize(t, f, 0)
	_, err := f.Write([]byte("surprise!"))
	if err == nil {
		checkSize(t, f, 13+9) // wrote at offset past where hello, world was.
	}
}

func TestTruncateNonexistentFile(t *testing.T) {
	t.Parallel()

	assertPathError := func(t testing.TB, path string, err error) {
		t.Helper()
		if pe, ok := err.(*PathError); !ok || !IsNotExist(err) || pe.Path != path {
			t.Errorf("got error: %v\nwant an ErrNotExist PathError with path %q", err, path)
		}
	}

	path := filepath.Join(t.TempDir(), "nonexistent")

	err := Truncate(path, 1)
	assertPathError(t, path, err)

	// Truncate shouldn't create any new file.
	_, err = Stat(path)
	assertPathError(t, path, err)
}

var hasNoatime = sync.OnceValue(func() bool {
	// A sloppy way to check if noatime flag is set (as all filesystems are
	// checked, not just the one we're interested in). A correct way
	// would be to use statvfs syscall and check if flags has ST_NOATIME,
	// but the syscall is OS-specific and is not even wired into Go stdlib.
	//
	// Only used on NetBSD (which ignores explicit atime updates with noatime).
	mounts, _ := ReadFile("/proc/mounts")
	return bytes.Contains(mounts, []byte("noatime"))
})

func TestChtimes(t *testing.T) {
	t.Parallel()

	f := newFile(t)
	// This should be an empty file (see #68687, #68663).
	f.Close()

	testChtimes(t, f.Name())
}

func TestChtimesOmit(t *testing.T) {
	t.Parallel()

	testChtimesOmit(t, true, false)
	testChtimesOmit(t, false, true)
	testChtimesOmit(t, true, true)
	testChtimesOmit(t, false, false) // Same as TestChtimes.
}

func testChtimesOmit(t *testing.T, omitAt, omitMt bool) {
	t.Logf("omit atime: %v, mtime: %v", omitAt, omitMt)
	file := newFile(t)
	// This should be an empty file (see #68687, #68663).
	name := file.Name()
	err := file.Close()
	if err != nil {
		t.Error(err)
	}
	fs, err := Stat(name)
	if err != nil {
		t.Fatal(err)
	}

	wantAtime := Atime(fs)
	wantMtime := fs.ModTime()
	switch runtime.GOOS {
	case "js":
		wantAtime = wantAtime.Truncate(time.Second)
		wantMtime = wantMtime.Truncate(time.Second)
	}

	var setAtime, setMtime time.Time // Zero value means omit.
	if !omitAt {
		wantAtime = wantAtime.Add(-1 * time.Second)
		setAtime = wantAtime
	}
	if !omitMt {
		wantMtime = wantMtime.Add(-1 * time.Second)
		setMtime = wantMtime
	}

	// Change the times accordingly.
	if err := Chtimes(name, setAtime, setMtime); err != nil {
		t.Error(err)
	}

	// Verify the expectations.
	fs, err = Stat(name)
	if err != nil {
		t.Error(err)
	}
	gotAtime := Atime(fs)
	gotMtime := fs.ModTime()

	// TODO: remove the dragonfly omitAt && omitMt exceptions below once the
	// fix (https://github.com/DragonFlyBSD/DragonFlyBSD/commit/c7c71870ed0)
	// is available generally and on CI runners.
	if !gotAtime.Equal(wantAtime) {
		errormsg := fmt.Sprintf("atime mismatch, got: %q, want: %q", gotAtime, wantAtime)
		switch runtime.GOOS {
		case "plan9", "tamago":
			// Mtime is the time of the last change of content.
			// Similarly, atime is set whenever the contents are
			// accessed; also, it is set whenever mtime is set.
		case "dragonfly":
			if omitAt && omitMt {
				t.Log(errormsg)
				t.Log("Known DragonFly BSD issue (won't work when both times are omitted); ignoring.")
			} else {
				// Assume hammer2 fs; https://www.dragonflybsd.org/hammer/ says:
				// > Because HAMMER2 is a block copy-on-write filesystem,
				// > the "atime" field is not supported and will typically
				// > just reflect local system in-memory caches or mtime.
				//
				// TODO: if only can CI define TMPDIR to point to a tmpfs
				// (e.g. /var/run/shm), this exception can be removed.
				t.Log(errormsg)
				t.Log("Known DragonFly BSD issue (atime not supported on hammer2); ignoring.")
			}
		case "netbsd":
			if !omitAt && hasNoatime() {
				t.Log(errormsg)
				t.Log("Known NetBSD issue (atime not changed on fs mounted with noatime); ignoring.")
			} else {
				t.Error(errormsg)
			}
		default:
			t.Error(errormsg)
		}
	}
	if !gotMtime.Equal(wantMtime) {
		errormsg := fmt.Sprintf("mtime mismatch, got: %q, want: %q", gotMtime, wantMtime)
		switch runtime.GOOS {
		case "dragonfly":
			if omitAt && omitMt {
				t.Log(errormsg)
				t.Log("Known DragonFly BSD issue (won't work when both times are omitted); ignoring.")
			} else {
				t.Error(errormsg)
			}
		default:
			t.Error(errormsg)
		}
	}
}

func TestChtimesDir(t *testing.T) {
	t.Parallel()

	testChtimes(t, t.TempDir())
}

func testChtimes(t *testing.T, name string) {
	st, err := Stat(name)
	if err != nil {
		t.Fatalf("Stat %s: %s", name, err)
	}
	preStat := st

	// Move access and modification time back a second
	at := Atime(preStat)
	mt := preStat.ModTime()
	err = Chtimes(name, at.Add(-time.Second), mt.Add(-time.Second))
	if err != nil {
		t.Fatalf("Chtimes %s: %s", name, err)
	}

	st, err = Stat(name)
	if err != nil {
		t.Fatalf("second Stat %s: %s", name, err)
	}
	postStat := st

	pat := Atime(postStat)
	pmt := postStat.ModTime()
	if !pat.Before(at) {
		errormsg := fmt.Sprintf("AccessTime didn't go backwards; was=%v, after=%v", at, pat)
		switch runtime.GOOS {
		case "plan9", "tamago":
			// Mtime is the time of the last change of
			// content.  Similarly, atime is set whenever
			// the contents are accessed; also, it is set
			// whenever mtime is set.
		case "netbsd":
			if hasNoatime() {
				t.Log(errormsg)
				t.Log("Known NetBSD issue (atime not changed on fs mounted with noatime); ignoring.")
			} else {
				t.Error(errormsg)
			}
		default:
			t.Error(errormsg)
		}
	}

	if !pmt.Before(mt) {
		t.Errorf("ModTime didn't go backwards; was=%v, after=%v", mt, pmt)
	}
}

func TestChtimesToUnixZero(t *testing.T) {
	file := newFile(t)
	fn := file.Name()
	if _, err := file.Write([]byte("hi")); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	unixZero := time.Unix(0, 0)
	if err := Chtimes(fn, unixZero, unixZero); err != nil {
		t.Fatalf("Chtimes failed: %v", err)
	}

	st, err := Stat(fn)
	if err != nil {
		t.Fatal(err)
	}

	if mt := st.ModTime(); mt != unixZero {
		t.Errorf("mtime is %v, want %v", mt, unixZero)
	}
}

func TestFileChdir(t *testing.T) {
	wd, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd: %s", err)
	}
	t.Chdir(".") // Ensure wd is restored after the test.

	fd, err := Open(".")
	if err != nil {
		t.Fatalf("Open .: %s", err)
	}
	defer fd.Close()

	if err := Chdir("/"); err != nil {
		t.Fatalf("Chdir /: %s", err)
	}

	if err := fd.Chdir(); err != nil {
		t.Fatalf("fd.Chdir: %s", err)
	}

	wdNew, err := Getwd()
	if err != nil {
		t.Fatalf("Getwd: %s", err)
	}

	wdInfo, err := fd.Stat()
	if err != nil {
		t.Fatal(err)
	}
	newInfo, err := Stat(wdNew)
	if err != nil {
		t.Fatal(err)
	}
	if !SameFile(wdInfo, newInfo) {
		t.Fatalf("fd.Chdir failed: got %s, want %s", wdNew, wd)
	}
}

func TestChdirAndGetwd(t *testing.T) {
	t.Chdir(t.TempDir()) // Ensure wd is restored after the test.

	// These are chosen carefully not to be symlinks on a Mac
	// (unlike, say, /var, /etc), except /tmp, which we handle below.
	dirs := []string{"/", "/usr/bin", "/tmp"}
	// /usr/bin does not usually exist on Plan 9 or Android.
	switch runtime.GOOS {
	case "android":
		dirs = []string{"/system/bin"}
	case "plan9":
		dirs = []string{"/", "/usr"}
	case "ios", "windows", "wasip1", "tamago":
		dirs = nil
		for _, dir := range []string{t.TempDir(), t.TempDir()} {
			// Expand symlinks so path equality tests work.
			dir, err := filepath.EvalSymlinks(dir)
			if err != nil {
				t.Fatalf("EvalSymlinks: %v", err)
			}
			dirs = append(dirs, dir)
		}
	}
	for mode := 0; mode < 2; mode++ {
		for _, d := range dirs {
			var err error
			if mode == 0 {
				err = Chdir(d)
			} else {
				fd1, err1 := Open(d)
				if err1 != nil {
					t.Errorf("Open %s: %s", d, err1)
					continue
				}
				err = fd1.Chdir()
				fd1.Close()
			}
			if d == "/tmp" {
				Setenv("PWD", "/tmp")
			}
			pwd, err1 := Getwd()
			if err != nil {
				t.Fatalf("Chdir %s: %s", d, err)
			}
			if err1 != nil {
				t.Fatalf("Getwd in %s: %s", d, err1)
			}
			if !equal(pwd, d) {
				t.Fatalf("Getwd returned %q want %q", pwd, d)
			}
		}
	}
}

// Test that Chdir+Getwd is program-wide.
func TestProgWideChdir(t *testing.T) {
	const N = 10
	var wg sync.WaitGroup
	hold := make(chan struct{})
	done := make(chan struct{})

	d := t.TempDir()
	t.Chdir(d)

	// Note the deferred Wait must be called after the deferred close(done),
	// to ensure the N goroutines have been released even if the main goroutine
	// calls Fatalf. It must be called before the Chdir back to the original
	// directory, and before the deferred deletion implied by TempDir,
	// so as not to interfere while the N goroutines are still running.
	defer wg.Wait()
	defer close(done)

	for i := 0; i < N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			// Lock half the goroutines in their own operating system
			// thread to exercise more scheduler possibilities.
			if i%2 == 1 {
				// On Plan 9, after calling LockOSThread, the goroutines
				// run on different processes which don't share the working
				// directory. This used to be an issue because Go expects
				// the working directory to be program-wide.
				// See issue 9428.
				runtime.LockOSThread()
			}
			select {
			case <-done:
				return
			case <-hold:
			}
			// Getwd might be wrong
			f0, err := Stat(".")
			if err != nil {
				t.Error(err)
				return
			}
			pwd, err := Getwd()
			if err != nil {
				t.Errorf("Getwd: %v", err)
				return
			}
			if pwd != d {
				t.Errorf("Getwd() = %q, want %q", pwd, d)
				return
			}
			f1, err := Stat(pwd)
			if err != nil {
				t.Error(err)
				return
			}
			if !SameFile(f0, f1) {
				t.Errorf(`Samefile(Stat("."), Getwd()) reports false (%s != %s)`, f0.Name(), f1.Name())
				return
			}
		}(i)
	}
	var err error
	if err = Chdir(d); err != nil {
		t.Fatalf("Chdir: %v", err)
	}
	// OS X sets TMPDIR to a symbolic link.
	// So we resolve our working directory again before the test.
	d, err = Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	close(hold)
	wg.Wait()
}

func TestSeek(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	const data = "hello, world\n"
	io.WriteString(f, data)

	type test struct {
		in     int64
		whence int
		out    int64
	}
	var tests = []test{
		{0, io.SeekCurrent, int64(len(data))},
		{0, io.SeekStart, 0},
		{5, io.SeekStart, 5},
		{0, io.SeekEnd, int64(len(data))},
		{0, io.SeekStart, 0},
		{-1, io.SeekEnd, int64(len(data)) - 1},
		{1 << 33, io.SeekStart, 1 << 33},
		{1 << 33, io.SeekEnd, 1<<33 + int64(len(data))},

		// Issue 21681, Windows 4G-1, etc:
		{1<<32 - 1, io.SeekStart, 1<<32 - 1},
		{0, io.SeekCurrent, 1<<32 - 1},
		{2<<32 - 1, io.SeekStart, 2<<32 - 1},
		{0, io.SeekCurrent, 2<<32 - 1},
	}
	for i, tt := range tests {
		off, err := f.Seek(tt.in, tt.whence)
		if off != tt.out || err != nil {
			t.Errorf("#%d: Seek(%v, %v) = %v, %v want %v, nil", i, tt.in, tt.whence, off, err, tt.out)
		}
	}
}

func TestSeekError(t *testing.T) {
	switch runtime.GOOS {
	case "js", "plan9", "wasip1", "tamago":
		t.Skipf("skipping test on %v", runtime.GOOS)
	}
	t.Parallel()

	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Seek(0, 0)
	if err == nil {
		t.Fatal("Seek on pipe should fail")
	}
	if perr, ok := err.(*PathError); !ok || perr.Err != syscall.ESPIPE {
		t.Errorf("Seek returned error %v, want &PathError{Err: syscall.ESPIPE}", err)
	}
	_, err = w.Seek(0, 0)
	if err == nil {
		t.Fatal("Seek on pipe should fail")
	}
	if perr, ok := err.(*PathError); !ok || perr.Err != syscall.ESPIPE {
		t.Errorf("Seek returned error %v, want &PathError{Err: syscall.ESPIPE}", err)
	}
}

func TestOpenError(t *testing.T) {
	t.Parallel()
	dir := makefs(t, []string{
		"is-a-file",
		"is-a-dir/",
	})
	t.Run("NoRoot", func(t *testing.T) { testOpenError(t, dir, false) })
	t.Run("InRoot", func(t *testing.T) { testOpenError(t, dir, true) })
}
func testOpenError(t *testing.T, dir string, rooted bool) {
	t.Parallel()
	var r *Root
	if rooted {
		var err error
		r, err = OpenRoot(dir)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
	}
	for _, tt := range []struct {
		path  string
		mode  int
		error error
	}{{
		"no-such-file",
		O_RDONLY,
		syscall.ENOENT,
	}, {
		"is-a-dir",
		O_WRONLY,
		syscall.EISDIR,
	}, {
		"is-a-file/no-such-file",
		O_WRONLY,
		syscall.ENOTDIR,
	}} {
		var f *File
		var err error
		var name string
		if rooted {
			name = fmt.Sprintf("Root(%q).OpenFile(%q, %d)", dir, tt.path, tt.mode)
			f, err = r.OpenFile(tt.path, tt.mode, 0)
		} else {
			path := filepath.Join(dir, tt.path)
			name = fmt.Sprintf("OpenFile(%q, %d)", path, tt.mode)
			f, err = OpenFile(path, tt.mode, 0)
		}
		if err == nil {
			t.Errorf("%v succeeded", name)
			f.Close()
			continue
		}
		perr, ok := err.(*PathError)
		if !ok {
			t.Errorf("%v returns error of %T type; want *PathError", name, err)
		}
		if perr.Err != tt.error {
			if runtime.GOOS == "plan9" {
				syscallErrStr := perr.Err.Error()
				expectedErrStr := strings.Replace(tt.error.Error(), "file ", "", 1)
				if !strings.HasSuffix(syscallErrStr, expectedErrStr) {
					// Some Plan 9 file servers incorrectly return
					// EPERM or EACCES rather than EISDIR when a directory is
					// opened for write.
					if tt.error == syscall.EISDIR &&
						(strings.HasSuffix(syscallErrStr, syscall.EPERM.Error()) ||
							strings.HasSuffix(syscallErrStr, syscall.EACCES.Error())) {
						continue
					}
					t.Errorf("%v = _, %q; want suffix %q", name, syscallErrStr, expectedErrStr)
				}
				continue
			}
			if runtime.GOOS == "dragonfly" {
				// DragonFly incorrectly returns EACCES rather
				// EISDIR when a directory is opened for write.
				if tt.error == syscall.EISDIR && perr.Err == syscall.EACCES {
					continue
				}
			}
			t.Errorf("%v = _, %q; want %q", name, perr.Err.Error(), tt.error.Error())
		}
	}
}

func TestOpenNoName(t *testing.T) {
	f, err := Open("")
	if err == nil {
		f.Close()
		t.Fatal(`Open("") succeeded`)
	}
}

func runBinHostname(t *testing.T) string {
	// Run /bin/hostname and collect output.
	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	path, err := exec.LookPath("hostname")
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			t.Skip("skipping test; test requires hostname but it does not exist")
		}
		t.Fatal(err)
	}

	argv := []string{"hostname"}
	if runtime.GOOS == "aix" {
		argv = []string{"hostname", "-s"}
	}
	p, err := StartProcess(path, argv, &ProcAttr{Files: []*File{nil, w, Stderr}})
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	var b strings.Builder
	io.Copy(&b, r)
	_, err = p.Wait()
	if err != nil {
		t.Fatalf("run hostname Wait: %v", err)
	}
	err = p.Kill()
	if err == nil {
		t.Errorf("expected an error from Kill running 'hostname'")
	}
	output := b.String()
	if n := len(output); n > 0 && output[n-1] == '\n' {
		output = output[0 : n-1]
	}
	if output == "" {
		t.Fatalf("/bin/hostname produced no output")
	}

	return output
}

func testWindowsHostname(t *testing.T, hostname string) {
	cmd := testenv.Command(t, "hostname")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to execute hostname command: %v %s", err, out)
	}
	want := strings.Trim(string(out), "\r\n")
	if hostname != want {
		t.Fatalf("Hostname() = %q != system hostname of %q", hostname, want)
	}
}

func TestHostname(t *testing.T) {
	t.Parallel()

	hostname, err := Hostname()
	if err != nil {
		t.Fatal(err)
	}
	if hostname == "" {
		t.Fatal("Hostname returned empty string and no error")
	}
	if strings.Contains(hostname, "\x00") {
		t.Fatalf("unexpected zero byte in hostname: %q", hostname)
	}

	// There is no other way to fetch hostname on windows, but via winapi.
	// On Plan 9 it can be taken from #c/sysname as Hostname() does.
	switch runtime.GOOS {
	case "android", "plan9":
		// No /bin/hostname to verify against.
		return
	case "windows":
		testWindowsHostname(t, hostname)
		return
	}

	testenv.MustHaveExec(t)

	// Check internal Hostname() against the output of /bin/hostname.
	// Allow that the internal Hostname returns a Fully Qualified Domain Name
	// and the /bin/hostname only returns the first component
	want := runBinHostname(t)
	if hostname != want {
		host, _, ok := strings.Cut(hostname, ".")
		if !ok || host != want {
			t.Errorf("Hostname() = %q, want %q", hostname, want)
		}
	}
}

func TestReadAt(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	const data = "hello, world\n"
	io.WriteString(f, data)

	b := make([]byte, 5)
	n, err := f.ReadAt(b, 7)
	if err != nil || n != len(b) {
		t.Fatalf("ReadAt 7: %d, %v", n, err)
	}
	if string(b) != "world" {
		t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
	}
}

// Verify that ReadAt doesn't affect seek offset.
// In the Plan 9 kernel, there used to be a bug in the implementation of
// the pread syscall, where the channel offset was erroneously updated after
// calling pread on a file.
func TestReadAtOffset(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	const data = "hello, world\n"
	io.WriteString(f, data)

	f.Seek(0, 0)
	b := make([]byte, 5)

	n, err := f.ReadAt(b, 7)
	if err != nil || n != len(b) {
		t.Fatalf("ReadAt 7: %d, %v", n, err)
	}
	if string(b) != "world" {
		t.Fatalf("ReadAt 7: have %q want %q", string(b), "world")
	}

	n, err = f.Read(b)
	if err != nil || n != len(b) {
		t.Fatalf("Read: %d, %v", n, err)
	}
	if string(b) != "hello" {
		t.Fatalf("Read: have %q want %q", string(b), "hello")
	}
}

// Verify that ReadAt doesn't allow negative offset.
func TestReadAtNegativeOffset(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	const data = "hello, world\n"
	io.WriteString(f, data)

	f.Seek(0, 0)
	b := make([]byte, 5)

	n, err := f.ReadAt(b, -10)

	const wantsub = "negative offset"
	if !strings.Contains(fmt.Sprint(err), wantsub) || n != 0 {
		t.Errorf("ReadAt(-10) = %v, %v; want 0, ...%q...", n, err, wantsub)
	}
}

func TestWriteAt(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	const data = "hello, world\n"
	io.WriteString(f, data)

	n, err := f.WriteAt([]byte("WORLD"), 7)
	if err != nil || n != 5 {
		t.Fatalf("WriteAt 7: %d, %v", n, err)
	}

	b, err := ReadFile(f.Name())
	if err != nil {
		t.Fatalf("ReadFile %s: %v", f.Name(), err)
	}
	if string(b) != "hello, WORLD\n" {
		t.Fatalf("after write: have %q want %q", string(b), "hello, WORLD\n")
	}
}

// Verify that WriteAt doesn't allow negative offset.
func TestWriteAtNegativeOffset(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	n, err := f.WriteAt([]byte("WORLD"), -10)

	const wantsub = "negative offset"
	if !strings.Contains(fmt.Sprint(err), wantsub) || n != 0 {
		t.Errorf("WriteAt(-10) = %v, %v; want 0, ...%q...", n, err, wantsub)
	}
}

// Verify that WriteAt doesn't work in append mode.
func TestWriteAtInAppendMode(t *testing.T) {
	t.Chdir(t.TempDir())
	f, err := OpenFile("write_at_in_append_mode.txt", O_APPEND|O_CREATE, 0666)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer f.Close()

	_, err = f.WriteAt([]byte(""), 1)
	if err != ErrWriteAtInAppendMode {
		t.Fatalf("f.WriteAt returned %v, expected %v", err, ErrWriteAtInAppendMode)
	}
}

func writeFile(t *testing.T, r *Root, fname string, flag int, text string) string {
	t.Helper()
	var f *File
	var err error
	if r == nil {
		f, err = OpenFile(fname, flag, 0666)
	} else {
		f, err = r.OpenFile(fname, flag, 0666)
	}
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	n, err := io.WriteString(f, text)
	if err != nil {
		t.Fatalf("WriteString: %d, %v", n, err)
	}
	f.Close()
	data, err := ReadFile(fname)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	return string(data)
}

func TestAppend(t *testing.T) {
	testMaybeRooted(t, func(t *testing.T, r *Root) {
		const f = "append.txt"
		s := writeFile(t, r, f, O_CREATE|O_TRUNC|O_RDWR, "new")
		if s != "new" {
			t.Fatalf("writeFile: have %q want %q", s, "new")
		}
		s = writeFile(t, r, f, O_APPEND|O_RDWR, "|append")
		if s != "new|append" {
			t.Fatalf("writeFile: have %q want %q", s, "new|append")
		}
		s = writeFile(t, r, f, O_CREATE|O_APPEND|O_RDWR, "|append")
		if s != "new|append|append" {
			t.Fatalf("writeFile: have %q want %q", s, "new|append|append")
		}
		err := Remove(f)
		if err != nil {
			t.Fatalf("Remove: %v", err)
		}
		s = writeFile(t, r, f, O_CREATE|O_APPEND|O_RDWR, "new&append")
		if s != "new&append" {
			t.Fatalf("writeFile: after append have %q want %q", s, "new&append")
		}
		s = writeFile(t, r, f, O_CREATE|O_RDWR, "old")
		if s != "old&append" {
			t.Fatalf("writeFile: after create have %q want %q", s, "old&append")
		}
		s = writeFile(t, r, f, O_CREATE|O_TRUNC|O_RDWR, "new")
		if s != "new" {
			t.Fatalf("writeFile: after truncate have %q want %q", s, "new")
		}
	})
}

// TestFilePermissions tests setting Unix permission bits on file creation.
func TestFilePermissions(t *testing.T) {
	if Getuid() == 0 {
		t.Skip("skipping test when running as root")
	}
	for _, test := range []struct {
		name string
		mode FileMode
	}{
		{"r", 0o444},
		{"w", 0o222},
		{"rw", 0o666},
	} {
		t.Run(test.name, func(t *testing.T) {
			switch runtime.GOOS {
			case "windows":
				if test.mode&0444 == 0 {
					t.Skip("write-only files not supported on " + runtime.GOOS)
				}
			case "wasip1", "tamago":
				t.Skip("file permissions not supported on " + runtime.GOOS)
			}
			testMaybeRooted(t, func(t *testing.T, r *Root) {
				const filename = "f"
				var f *File
				var err error
				if r == nil {
					f, err = OpenFile(filename, O_RDWR|O_CREATE|O_EXCL, test.mode)
				} else {
					f, err = r.OpenFile(filename, O_RDWR|O_CREATE|O_EXCL, test.mode)
				}
				if err != nil {
					t.Fatal(err)
				}
				f.Close()
				b, err := ReadFile(filename)
				if test.mode&0o444 != 0 {
					if err != nil {
						t.Errorf("ReadFile = %v; want success", err)
					}
				} else {
					if err == nil {
						t.Errorf("ReadFile = %q, <nil>; want failure", string(b))
					}
				}
				_, err = Stat(filename)
				if err != nil {
					t.Errorf("Stat = %v; want success", err)
				}
				err = WriteFile(filename, nil, 0666)
				if test.mode&0o222 != 0 {
					if err != nil {
						t.Errorf("WriteFile = %v; want success", err)
						b, err := ReadFile(filename)
						t.Errorf("ReadFile: %v", err)
						t.Errorf("file contents: %q", b)
					}
				} else {
					if err == nil {
						t.Errorf("WriteFile(%q) = <nil>; want failure", filename)
						st, err := Stat(filename)
						if err == nil {
							t.Errorf("mode: %s", st.Mode())
						}
						b, err := ReadFile(filename)
						t.Errorf("ReadFile: %v", err)
						t.Errorf("file contents: %q", b)
					}
				}
			})
		})
	}

}

func TestOpenFileCreateExclDanglingSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)

	testMaybeRooted(t, func(t *testing.T, r *Root) {
		const link = "link"
		if err := Symlink("does_not_exist", link); err != nil {
			t.Fatal(err)
		}
		var f *File
		var err error
		if r == nil {
			f, err = OpenFile(link, O_WRONLY|O_CREATE|O_EXCL, 0o444)
		} else {
			f, err = r.OpenFile(link, O_WRONLY|O_CREATE|O_EXCL, 0o444)
		}
		if err == nil {
			f.Close()
		}
		if !errors.Is(err, ErrExist) {
			t.Errorf("OpenFile of a dangling symlink with O_CREATE|O_EXCL = %v, want ErrExist", err)
		}
		if _, err := Stat(link); err == nil {
			t.Errorf("OpenFile of a dangling symlink with O_CREATE|O_EXCL created a file")
		}
	})
}

// TestFileRDWRFlags tests the O_RDONLY, O_WRONLY, and O_RDWR flags.
func TestFileRDWRFlags(t *testing.T) {
	for _, test := range []struct {
		name string
		flag int
	}{
		{"O_RDONLY", O_RDONLY},
		{"O_WRONLY", O_WRONLY},
		{"O_RDWR", O_RDWR},
	} {
		t.Run(test.name, func(t *testing.T) {
			testMaybeRooted(t, func(t *testing.T, r *Root) {
				const filename = "f"
				content := []byte("content")
				if err := WriteFile(filename, content, 0666); err != nil {
					t.Fatal(err)
				}
				var f *File
				var err error
				if r == nil {
					f, err = OpenFile(filename, test.flag, 0)
				} else {
					f, err = r.OpenFile(filename, test.flag, 0)
				}
				if err != nil {
					t.Fatal(err)
				}
				defer f.Close()
				got, err := io.ReadAll(f)
				if test.flag == O_WRONLY {
					if err == nil {
						t.Errorf("read file: %q, %v; want error", got, err)
					}
				} else {
					if err != nil || !bytes.Equal(got, content) {
						t.Errorf("read file: %q, %v; want %q, <nil>", got, err, content)
					}
				}
				if _, err := f.Seek(0, 0); err != nil {
					t.Fatalf("f.Seek: %v", err)
				}
				newcontent := []byte("CONTENT")
				_, err = f.Write(newcontent)
				if test.flag == O_RDONLY {
					if err == nil {
						t.Errorf("write file: succeeded, want error")
					}
				} else {
					if err != nil {
						t.Errorf("write file: %v, want success", err)
					}
				}
				f.Close()
				got, err = ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				want := content
				if test.flag != O_RDONLY {
					want = newcontent
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("after write, file contains %q, want %q", got, want)
				}
			})
		})
	}
}

func TestStatDirWithTrailingSlash(t *testing.T) {
	t.Parallel()

	// Create new temporary directory and arrange to clean it up.
	path := t.TempDir()

	// Stat of path should succeed.
	if _, err := Stat(path); err != nil {
		t.Fatalf("stat %s failed: %s", path, err)
	}

	// Stat of path+"/" should succeed too.
	path += "/"
	if _, err := Stat(path); err != nil {
		t.Fatalf("stat %s failed: %s", path, err)
	}
}

func TestNilProcessStateString(t *testing.T) {
	var ps *ProcessState
	s := ps.String()
	if s != "<nil>" {
		t.Errorf("(*ProcessState)(nil).String() = %q, want %q", s, "<nil>")
	}
}

func TestSameFile(t *testing.T) {
	t.Chdir(t.TempDir())
	fa, err := Create("a")
	if err != nil {
		t.Fatalf("Create(a): %v", err)
	}
	fa.Close()
	fb, err := Create("b")
	if err != nil {
		t.Fatalf("Create(b): %v", err)
	}
	fb.Close()

	ia1, err := Stat("a")
	if err != nil {
		t.Fatalf("Stat(a): %v", err)
	}
	ia2, err := Stat("a")
	if err != nil {
		t.Fatalf("Stat(a): %v", err)
	}
	if !SameFile(ia1, ia2) {
		t.Errorf("files should be same")
	}

	ib, err := Stat("b")
	if err != nil {
		t.Fatalf("Stat(b): %v", err)
	}
	if SameFile(ia1, ib) {
		t.Errorf("files should be different")
	}
}

func testDevNullFileInfo(t *testing.T, statname, devNullName string, fi FileInfo) {
	pre := fmt.Sprintf("%s(%q): ", statname, devNullName)
	if fi.Size() != 0 {
		t.Errorf(pre+"wrong file size have %d want 0", fi.Size())
	}
	if fi.Mode()&ModeDevice == 0 {
		t.Errorf(pre+"wrong file mode %q: ModeDevice is not set", fi.Mode())
	}
	if fi.Mode()&ModeCharDevice == 0 {
		t.Errorf(pre+"wrong file mode %q: ModeCharDevice is not set", fi.Mode())
	}
	if fi.Mode().IsRegular() {
		t.Errorf(pre+"wrong file mode %q: IsRegular returns true", fi.Mode())
	}
}

func testDevNullFile(t *testing.T, devNullName string) {
	f, err := Open(devNullName)
	if err != nil {
		t.Fatalf("Open(%s): %v", devNullName, err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		t.Fatalf("Stat(%s): %v", devNullName, err)
	}
	testDevNullFileInfo(t, "f.Stat", devNullName, fi)

	fi, err = Stat(devNullName)
	if err != nil {
		t.Fatalf("Stat(%s): %v", devNullName, err)
	}
	testDevNullFileInfo(t, "Stat", devNullName, fi)
}

func TestDevNullFile(t *testing.T) {
	t.Parallel()

	testDevNullFile(t, DevNull)
	if runtime.GOOS == "windows" {
		testDevNullFile(t, "./nul")
		testDevNullFile(t, "//./nul")
	}
}

var testLargeWrite = flag.Bool("large_write", false, "run TestLargeWriteToConsole test that floods console with output")

func TestLargeWriteToConsole(t *testing.T) {
	if !*testLargeWrite {
		t.Skip("skipping console-flooding test; enable with -large_write")
	}
	b := make([]byte, 32000)
	for i := range b {
		b[i] = '.'
	}
	b[len(b)-1] = '\n'
	n, err := Stdout.Write(b)
	if err != nil {
		t.Fatalf("Write to os.Stdout failed: %v", err)
	}
	if n != len(b) {
		t.Errorf("Write to os.Stdout should return %d; got %d", len(b), n)
	}
	n, err = Stderr.Write(b)
	if err != nil {
		t.Fatalf("Write to os.Stderr failed: %v", err)
	}
	if n != len(b) {
		t.Errorf("Write to os.Stderr should return %d; got %d", len(b), n)
	}
}

func TestStatDirModeExec(t *testing.T) {
	if runtime.GOOS == "wasip1" {
		t.Skip("Chmod is not supported on " + runtime.GOOS)
	}
	t.Parallel()

	const mode = 0111

	path := t.TempDir()
	if err := Chmod(path, 0777); err != nil {
		t.Fatalf("Chmod %q 0777: %v", path, err)
	}

	dir, err := Stat(path)
	if err != nil {
		t.Fatalf("Stat %q (looking for mode %#o): %s", path, mode, err)
	}
	if dir.Mode()&mode != mode {
		t.Errorf("Stat %q: mode %#o want %#o", path, dir.Mode()&mode, mode)
	}
}

func TestStatStdin(t *testing.T) {
	switch runtime.GOOS {
	case "android", "plan9":
		t.Skipf("%s doesn't have /bin/sh", runtime.GOOS)
	}

	if Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		st, err := Stdin.Stat()
		if err != nil {
			t.Fatalf("Stat failed: %v", err)
		}
		fmt.Println(st.Mode() & ModeNamedPipe)
		Exit(0)
	}

	t.Parallel()
	exe := testenv.Executable(t)

	fi, err := Stdin.Stat()
	if err != nil {
		t.Fatal(err)
	}
	switch mode := fi.Mode(); {
	case mode&ModeCharDevice != 0 && mode&ModeDevice != 0:
	case mode&ModeNamedPipe != 0:
	default:
		t.Fatalf("unexpected Stdin mode (%v), want ModeCharDevice or ModeNamedPipe", mode)
	}

	cmd := testenv.Command(t, exe, "-test.run=^TestStatStdin$")
	cmd = testenv.CleanCmdEnv(cmd)
	cmd.Env = append(cmd.Env, "GO_WANT_HELPER_PROCESS=1")
	// This will make standard input a pipe.
	cmd.Stdin = strings.NewReader("output")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to spawn child process: %v %q", err, string(output))
	}

	// result will be like "prw-rw-rw"
	if len(output) < 1 || output[0] != 'p' {
		t.Fatalf("Child process reports stdin is not pipe '%v'", string(output))
	}
}

func TestStatRelativeSymlink(t *testing.T) {
	testenv.MustHaveSymlink(t)
	t.Parallel()

	tmpdir := t.TempDir()
	target := filepath.Join(tmpdir, "target")
	f, err := Create(target)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(tmpdir, "link")
	err = Symlink(filepath.Base(target), link)
	if err != nil {
		t.Fatal(err)
	}

	st1, err := Stat(link)
	if err != nil {
		t.Fatal(err)
	}

	if !SameFile(st, st1) {
		t.Error("Stat doesn't follow relative symlink")
	}

	if runtime.GOOS == "windows" {
		Remove(link)
		err = Symlink(target[len(filepath.VolumeName(target)):], link)
		if err != nil {
			t.Fatal(err)
		}

		st1, err := Stat(link)
		if err != nil {
			t.Fatal(err)
		}

		if !SameFile(st, st1) {
			t.Error("Stat doesn't follow relative symlink")
		}
	}
}

func TestReadAtEOF(t *testing.T) {
	t.Parallel()

	f := newFile(t)

	_, err := f.ReadAt(make([]byte, 10), 0)
	switch err {
	case io.EOF:
		// all good
	case nil:
		t.Fatalf("ReadAt succeeded")
	default:
		t.Fatalf("ReadAt failed: %s", err)
	}
}

func TestLongPath(t *testing.T) {
	t.Parallel()

	tmpdir := t.TempDir()

	// Test the boundary of 247 and fewer bytes (normal) and 248 and more bytes (adjusted).
	sizes := []int{247, 248, 249, 400}
	for len(tmpdir) < 400 {
		tmpdir += "/dir3456789"
	}
	for _, sz := range sizes {
		t.Run(fmt.Sprintf("length=%d", sz), func(t *testing.T) {
			sizedTempDir := tmpdir[:sz-1] + "x" // Ensure it does not end with a slash.

			// The various sized runs are for this call to trigger the boundary
			// condition.
			if err := MkdirAll(sizedTempDir, 0755); err != nil {
				t.Fatalf("MkdirAll failed: %v", err)
			}
			data := []byte("hello world\n")
			if err := WriteFile(sizedTempDir+"/foo.txt", data, 0644); err != nil {
				t.Fatalf("os.WriteFile() failed: %v", err)
			}
			if err := Rename(sizedTempDir+"/foo.txt", sizedTempDir+"/bar.txt"); err != nil {
				t.Fatalf("Rename failed: %v", err)
			}
			mtime := time.Now().Truncate(time.Minute)
			if err := Chtimes(sizedTempDir+"/bar.txt", mtime, mtime); err != nil {
				t.Fatalf("Chtimes failed: %v", err)
			}
			names := []string{"bar.txt"}
			if testenv.HasSymlink() {
				if err := Symlink(sizedTempDir+"/bar.txt", sizedTempDir+"/symlink.txt"); err != nil {
					t.Fatalf("Symlink failed: %v", err)
				}
				names = append(names, "symlink.txt")
			}
			if testenv.HasLink() {
				if err := Link(sizedTempDir+"/bar.txt", sizedTempDir+"/link.txt"); err != nil {
					t.Fatalf("Link failed: %v", err)
				}
				names = append(names, "link.txt")
			}
			for _, wantSize := range []int64{int64(len(data)), 0} {
				for _, name := range names {
					path := sizedTempDir + "/" + name
					dir, err := Stat(path)
					if err != nil {
						t.Fatalf("Stat(%q) failed: %v", path, err)
					}
					filesize := size(path, t)
					if dir.Size() != filesize || filesize != wantSize {
						t.Errorf("Size(%q) is %d, len(ReadFile()) is %d, want %d", path, dir.Size(), filesize, wantSize)
					}
					if runtime.GOOS != "wasip1" { // Chmod is not supported on wasip1
						err = Chmod(path, dir.Mode())
						if err != nil {
							t.Fatalf("Chmod(%q) failed: %v", path, err)
						}
					}
				}
				if err := Truncate(sizedTempDir+"/bar.txt", 0); err != nil {
					t.Fatalf("Truncate failed: %v", err)
				}
			}
		})
	}
}

func testKillProcess(t *testing.T, processKiller func(p *Process)) {
	t.Parallel()

	// Re-exec the test binary to start a process that hangs until stdin is closed.
	cmd := testenv.Command(t, testenv.Executable(t))
	cmd.Env = append(cmd.Environ(), "GO_OS_TEST_DRAIN_STDIN=1")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start test process: %v", err)
	}

	defer func() {
		if err := cmd.Wait(); err == nil {
			t.Errorf("Test process succeeded, but expected to fail")
		}
		stdin.Close() // Keep stdin alive until the process has finished dying.
	}()

	// Wait for the process to be started.
	// (It will close its stdout when it reaches TestMain.)
	io.Copy(io.Discard, stdout)

	processKiller(cmd.Process)
}

func TestKillStartProcess(t *testing.T) {
	testKillProcess(t, func(p *Process) {
		err := p.Kill()
		if err != nil {
			t.Fatalf("Failed to kill test process: %v", err)
		}
	})
}

func TestGetppid(t *testing.T) {
	if runtime.GOOS == "plan9" {
		// TODO: golang.org/issue/8206
		t.Skipf("skipping test on plan9; see issue 8206")
	}

	if Getenv("GO_WANT_HELPER_PROCESS") == "1" {
		fmt.Print(Getppid())
		Exit(0)
	}

	t.Parallel()

	cmd := testenv.Command(t, testenv.Executable(t), "-test.run=^TestGetppid$")
	cmd.Env = append(Environ(), "GO_WANT_HELPER_PROCESS=1")

	// verify that Getppid() from the forked process reports our process id
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to spawn child process: %v %q", err, string(output))
	}

	childPpid := string(output)
	ourPid := fmt.Sprintf("%d", Getpid())
	if childPpid != ourPid {
		t.Fatalf("Child process reports parent process id '%v', expected '%v'", childPpid, ourPid)
	}
}

func TestKillFindProcess(t *testing.T) {
	testKillProcess(t, func(p *Process) {
		p2, err := FindProcess(p.Pid)
		if err != nil {
			t.Fatalf("Failed to find test process: %v", err)
		}
		err = p2.Kill()
		if err != nil {
			t.Fatalf("Failed to kill test process: %v", err)
		}
	})
}

var nilFileMethodTests = []struct {
	name string
	f    func(*File) error
}{
	{"Chdir", func(f *File) error { return f.Chdir() }},
	{"Close", func(f *File) error { return f.Close() }},
	{"Chmod", func(f *File) error { return f.Chmod(0) }},
	{"Chown", func(f *File) error { return f.Chown(0, 0) }},
	{"Read", func(f *File) error { _, err := f.Read(make([]byte, 0)); return err }},
	{"ReadAt", func(f *File) error { _, err := f.ReadAt(make([]byte, 0), 0); return err }},
	{"Readdir", func(f *File) error { _, err := f.Readdir(1); return err }},
	{"Readdirnames", func(f *File) error { _, err := f.Readdirnames(1); return err }},
	{"Seek", func(f *File) error { _, err := f.Seek(0, io.SeekStart); return err }},
	{"Stat", func(f *File) error { _, err := f.Stat(); return err }},
	{"Sync", func(f *File) error { return f.Sync() }},
	{"Truncate", func(f *File) error { return f.Truncate(0) }},
	{"Write", func(f *File) error { _, err := f.Write(make([]byte, 0)); return err }},
	{"WriteAt", func(f *File) error { _, err := f.WriteAt(make([]byte, 0), 0); return err }},
	{"WriteString", func(f *File) error { _, err := f.WriteString(""); return err }},
}

// Test that all File methods give ErrInvalid if the receiver is nil.
func TestNilFileMethods(t *testing.T) {
	t.Parallel()

	for _, tt := range nilFileMethodTests {
		var file *File
		got := tt.f(file)
		if got != ErrInvalid {
			t.Errorf("%v should fail when f is nil; got %v", tt.name, got)
		}
	}
}

func mkdirTree(t *testing.T, root string, level, max int) {
	if level >= max {
		return
	}
	level++
	for i := 'a'; i < 'c'; i++ {
		dir := filepath.Join(root, string(i))
		if err := Mkdir(dir, 0700); err != nil {
			t.Fatal(err)
		}
		mkdirTree(t, dir, level, max)
	}
}

// Test that simultaneous RemoveAll do not report an error.
// As long as it gets removed, we should be happy.
func TestRemoveAllRace(t *testing.T) {
	if runtime.GOOS == "windows" {
		// Windows has very strict rules about things like
		// removing directories while someone else has
		// them open. The racing doesn't work out nicely
		// like it does on Unix.
		t.Skip("skipping on windows")
	}
	if runtime.GOOS == "dragonfly" {
		testenv.SkipFlaky(t, 52301)
	}

	n := runtime.GOMAXPROCS(16)
	defer runtime.GOMAXPROCS(n)
	root := t.TempDir()
	mkdirTree(t, root, 1, 6)
	hold := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-hold
			err := RemoveAll(root)
			if err != nil {
				t.Errorf("unexpected error: %T, %q", err, err)
			}
		}()
	}
	close(hold) // let workers race to remove root
	wg.Wait()
}

// Test that reading from a pipe doesn't use up a thread.
func TestPipeThreads(t *testing.T) {
	switch runtime.GOOS {
	case "aix":
		t.Skip("skipping on aix; issue 70131")
	case "illumos", "solaris":
		t.Skip("skipping on Solaris and illumos; issue 19111")
	case "windows":
		t.Skip("skipping on Windows; issue 19098")
	case "plan9":
		t.Skip("skipping on Plan 9; does not support runtime poller")
	case "js":
		t.Skip("skipping on js; no support for os.Pipe")
	case "wasip1":
		t.Skip("skipping on wasip1; no support for os.Pipe")
	case "tamago":
		t.Skip("skipping on tamago; no support for os.Pipe")
	}

	threads := 100

	r := make([]*File, threads)
	w := make([]*File, threads)
	for i := 0; i < threads; i++ {
		rp, wp, err := Pipe()
		if err != nil {
			for j := 0; j < i; j++ {
				r[j].Close()
				w[j].Close()
			}
			t.Fatal(err)
		}
		r[i] = rp
		w[i] = wp
	}

	defer debug.SetMaxThreads(debug.SetMaxThreads(threads / 2))

	creading := make(chan bool, threads)
	cdone := make(chan bool, threads)
	for i := 0; i < threads; i++ {
		go func(i int) {
			var b [1]byte
			creading <- true
			if _, err := r[i].Read(b[:]); err != nil {
				t.Error(err)
			}
			if err := r[i].Close(); err != nil {
				t.Error(err)
			}
			cdone <- true
		}(i)
	}

	for i := 0; i < threads; i++ {
		<-creading
	}

	// If we are still alive, it means that the 100 goroutines did
	// not require 100 threads.

	for i := 0; i < threads; i++ {
		if _, err := w[i].Write([]byte{0}); err != nil {
			t.Error(err)
		}
		if err := w[i].Close(); err != nil {
			t.Error(err)
		}
		<-cdone
	}
}

func testDoubleCloseError(path string) func(*testing.T) {
	return func(t *testing.T) {
		t.Parallel()

		file, err := Open(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatalf("unexpected error from Close: %v", err)
		}
		if err := file.Close(); err == nil {
			t.Error("second Close did not fail")
		} else if pe, ok := err.(*PathError); !ok {
			t.Errorf("second Close: got %T, want %T", err, pe)
		} else if pe.Err != ErrClosed {
			t.Errorf("second Close: got %q, want %q", pe.Err, ErrClosed)
		} else {
			t.Logf("second close returned expected error %q", err)
		}
	}
}

func TestDoubleCloseError(t *testing.T) {
	t.Parallel()
	t.Run("file", testDoubleCloseError(filepath.Join(sfdir, sfname)))
	t.Run("dir", testDoubleCloseError(sfdir))
}

func TestUserCacheDir(t *testing.T) {
	t.Parallel()

	dir, err := UserCacheDir()
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	if dir == "" {
		t.Fatalf("UserCacheDir returned %q; want non-empty path or error", dir)
	}

	fi, err := Stat(dir)
	if err != nil {
		if IsNotExist(err) {
			t.Log(err)
			return
		}
		t.Fatal(err)
	}
	if !fi.IsDir() {
		t.Fatalf("dir %s is not directory; type = %v", dir, fi.Mode())
	}
}

func TestUserCacheDirXDGConfigDirEnvVar(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "darwin", "plan9", "tamago":
		t.Skip("$XDG_CACHE_HOME is effective only on Unix systems")
	}

	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CACHE_HOME", wd)

	dir, err := UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != wd {
		t.Fatalf("UserCacheDir returned %q; want the value of $XDG_CACHE_HOME %q", dir, wd)
	}

	t.Setenv("XDG_CACHE_HOME", "some-dir")
	_, err = UserCacheDir()
	if err == nil {
		t.Fatal("UserCacheDir succeeded though $XDG_CACHE_HOME contains a relative path")
	}
}

func TestUserConfigDir(t *testing.T) {
	t.Parallel()

	dir, err := UserConfigDir()
	if err != nil {
		t.Skipf("skipping: %v", err)
	}
	if dir == "" {
		t.Fatalf("UserConfigDir returned %q; want non-empty path or error", dir)
	}

	fi, err := Stat(dir)
	if err != nil {
		if IsNotExist(err) {
			t.Log(err)
			return
		}
		t.Fatal(err)
	}
	if !fi.IsDir() {
		t.Fatalf("dir %s is not directory; type = %v", dir, fi.Mode())
	}
}

func TestUserConfigDirXDGConfigDirEnvVar(t *testing.T) {
	switch runtime.GOOS {
	case "windows", "darwin", "plan9", "tamago":
		t.Skip("$XDG_CONFIG_HOME is effective only on Unix systems")
	}

	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", wd)

	dir, err := UserConfigDir()
	if err != nil {
		t.Fatal(err)
	}
	if dir != wd {
		t.Fatalf("UserConfigDir returned %q; want the value of $XDG_CONFIG_HOME %q", dir, wd)
	}

	t.Setenv("XDG_CONFIG_HOME", "some-dir")
	_, err = UserConfigDir()
	if err == nil {
		t.Fatal("UserConfigDir succeeded though $XDG_CONFIG_HOME contains a relative path")
	}
}

func TestUserHomeDir(t *testing.T) {
	t.Parallel()

	dir, err := UserHomeDir()
	if dir == "" && err == nil {
		t.Fatal("UserHomeDir returned an empty string but no error")
	}
	if err != nil {
		// UserHomeDir may return a non-nil error if the environment variable
		// for the home directory is empty or unset in the environment.
		t.Skipf("skipping: %v", err)
	}

	fi, err := Stat(dir)
	if err != nil {
		if IsNotExist(err) {
			// The user's home directory has a well-defined location, but does not
			// exist. (Maybe nothing has written to it yet? That could happen, for
			// example, on minimal VM images used for CI testing.)
			t.Log(err)
			return
		}
		t.Fatal(err)
	}
	if !fi.IsDir() {
		t.Fatalf("dir %s is not directory; type = %v", dir, fi.Mode())
	}
}

func TestDirSeek(t *testing.T) {
	t.Parallel()

	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	f, err := Open(wd)
	if err != nil {
		t.Fatal(err)
	}
	dirnames1, err := f.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	ret, err := f.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if ret != 0 {
		t.Fatalf("seek result not zero: %d", ret)
	}

	dirnames2, err := f.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}

	if len(dirnames1) != len(dirnames2) {
		t.Fatalf("listings have different lengths: %d and %d\n", len(dirnames1), len(dirnames2))
	}
	for i, n1 := range dirnames1 {
		n2 := dirnames2[i]
		if n1 != n2 {
			t.Fatalf("different name i=%d n1=%s n2=%s\n", i, n1, n2)
		}
	}
}

func TestReaddirSmallSeek(t *testing.T) {
	// See issue 37161. Read only one entry from a directory,
	// seek to the beginning, and read again. We should not see
	// duplicate entries.
	t.Parallel()

	wd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	df, err := Open(filepath.Join(wd, "testdata", "issue37161"))
	if err != nil {
		t.Fatal(err)
	}
	names1, err := df.Readdirnames(1)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = df.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	names2, err := df.Readdirnames(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(names2) != 3 {
		t.Fatalf("first names: %v, second names: %v", names1, names2)
	}
}

// isDeadlineExceeded reports whether err is or wraps ErrDeadlineExceeded.
// We also check that the error has a Timeout method that returns true.
func isDeadlineExceeded(err error) bool {
	if !IsTimeout(err) {
		return false
	}
	if !errors.Is(err, ErrDeadlineExceeded) {
		return false
	}
	return true
}

// Test that opening a file does not change its permissions.  Issue 38225.
func TestOpenFileKeepsPermissions(t *testing.T) {
	t.Run("OpenFile", func(t *testing.T) {
		testOpenFileKeepsPermissions(t, OpenFile)
	})
	t.Run("RootOpenFile", func(t *testing.T) {
		testOpenFileKeepsPermissions(t, func(name string, flag int, perm FileMode) (*File, error) {
			dir, file := filepath.Split(name)
			r, err := OpenRoot(dir)
			if err != nil {
				return nil, err
			}
			defer r.Close()
			return r.OpenFile(file, flag, perm)
		})
	})
}
func testOpenFileKeepsPermissions(t *testing.T, openf func(name string, flag int, perm FileMode) (*File, error)) {
	t.Parallel()

	dir := t.TempDir()
	name := filepath.Join(dir, "x")
	f, err := Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}
	f, err = openf(name, O_WRONLY|O_CREATE|O_TRUNC, 0)
	if err != nil {
		t.Fatal(err)
	}
	if fi, err := f.Stat(); err != nil {
		t.Error(err)
	} else if fi.Mode()&0222 == 0 {
		t.Errorf("f.Stat.Mode after OpenFile is %v, should be writable", fi.Mode())
	}
	if err := f.Close(); err != nil {
		t.Error(err)
	}
	if fi, err := Stat(name); err != nil {
		t.Error(err)
	} else if fi.Mode()&0222 == 0 {
		t.Errorf("Stat after OpenFile is %v, should be writable", fi.Mode())
	}
}

func forceMFTUpdateOnWindows(t *testing.T, path string) {
	t.Helper()

	if runtime.GOOS != "windows" {
		return
	}

	// On Windows, we force the MFT to update by reading the actual metadata from GetFileInformationByHandle and then
	// explicitly setting that. Otherwise it might get out of sync with FindFirstFile. See golang.org/issues/42637.
	if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			t.Fatal(err)
		}
		info, err := d.Info()
		if err != nil {
			t.Fatal(err)
		}
		stat, err := Stat(path) // This uses GetFileInformationByHandle internally.
		if err != nil {
			t.Fatal(err)
		}
		if stat.ModTime() == info.ModTime() {
			return nil
		}
		if err := Chtimes(path, stat.ModTime(), stat.ModTime()); err != nil {
			t.Log(err) // We only log, not die, in case the test directory is not writable.
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDirFS(t *testing.T) {
	t.Parallel()
	testDirFS(t, DirFS("./testdata/dirfs"))
}

func TestRootDirFS(t *testing.T) {
	t.Parallel()
	r, err := OpenRoot("./testdata/dirfs")
	if err != nil {
		t.Fatal(err)
	}
	testDirFS(t, r.FS())
}

func testDirFS(t *testing.T, fsys fs.FS) {
	forceMFTUpdateOnWindows(t, "./testdata/dirfs")

	if err := fstest.TestFS(fsys, "a", "b", "dir/x"); err != nil {
		t.Fatal(err)
	}

	rdfs, ok := fsys.(fs.ReadDirFS)
	if !ok {
		t.Error("expected DirFS result to implement fs.ReadDirFS")
	}
	if _, err := rdfs.ReadDir("nonexistent"); err == nil {
		t.Error("fs.ReadDir of nonexistent directory succeeded")
	}

	// Test that the error message does not contain a backslash,
	// and does not contain the DirFS argument.
	const nonesuch = "dir/nonesuch"
	_, err := fsys.Open(nonesuch)
	if err == nil {
		t.Error("fs.Open of nonexistent file succeeded")
	} else {
		if !strings.Contains(err.Error(), nonesuch) {
			t.Errorf("error %q does not contain %q", err, nonesuch)
		}
		if strings.Contains(err.(*PathError).Path, "testdata") {
			t.Errorf("error %q contains %q", err, "testdata")
		}
	}

	// Test that Open does not accept backslash as separator.
	d := DirFS(".")
	_, err = d.Open(`testdata\dirfs`)
	if err == nil {
		t.Fatalf(`Open testdata\dirfs succeeded`)
	}

	// Test that Open does not open Windows device files.
	_, err = d.Open(`NUL`)
	if err == nil {
		t.Errorf(`Open NUL succeeded`)
	}
}

func TestDirFSRootDir(t *testing.T) {
	t.Parallel()

	cwd, err := Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cwd = cwd[len(filepath.VolumeName(cwd)):] // trim volume prefix (C:) on Windows
	cwd = filepath.ToSlash(cwd)               // convert \ to /
	cwd = strings.TrimPrefix(cwd, "/")        // trim leading /

	// Test that Open can open a path starting at /.
	d := DirFS("/")
	f, err := d.Open(cwd + "/testdata/dirfs/a")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestDirFSEmptyDir(t *testing.T) {
	t.Parallel()

	d := DirFS("")
	cwd, _ := Getwd()
	for _, path := range []string{
		"testdata/dirfs/a",                          // not DirFS(".")
		filepath.ToSlash(cwd) + "/testdata/dirfs/a", // not DirFS("/")
	} {
		_, err := d.Open(path)
		if err == nil {
			t.Fatalf(`DirFS("").Open(%q) succeeded`, path)
		}
	}
}

func TestDirFSPathsValid(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("skipping on Windows")
	}
	t.Parallel()

	d := t.TempDir()
	if err := WriteFile(filepath.Join(d, "control.txt"), []byte(string("Hello, world!")), 0644); err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(filepath.Join(d, `e:xperi\ment.txt`), []byte(string("Hello, colon and backslash!")), 0644); err != nil {
		t.Fatal(err)
	}

	fsys := DirFS(d)
	err := fs.WalkDir(fsys, ".", func(path string, e fs.DirEntry, err error) error {
		if fs.ValidPath(e.Name()) {
			t.Logf("%q ok", e.Name())
		} else {
			t.Errorf("%q INVALID", e.Name())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadFileProc(t *testing.T) {
	t.Parallel()

	// Linux files in /proc report 0 size,
	// but then if ReadFile reads just a single byte at offset 0,
	// the read at offset 1 returns EOF instead of more data.
	// ReadFile has a minimum read size of 512 to work around this,
	// but test explicitly that it's working.
	name := "/proc/sys/fs/pipe-max-size"
	if _, err := Stat(name); err != nil {
		t.Skip(err)
	}
	data, err := ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatalf("read %s: not newline-terminated: %q", name, data)
	}
}

func TestDirFSReadFileProc(t *testing.T) {
	t.Parallel()

	fsys := DirFS("/")
	name := "proc/sys/fs/pipe-max-size"
	if _, err := fs.Stat(fsys, name); err != nil {
		t.Skip()
	}
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || data[len(data)-1] != '\n' {
		t.Fatalf("read %s: not newline-terminated: %q", name, data)
	}
}

func TestWriteStringAlloc(t *testing.T) {
	if runtime.GOOS == "js" {
		t.Skip("js allocates a lot during File.WriteString")
	}
	d := t.TempDir()
	f, err := Create(filepath.Join(d, "whiteboard.txt"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	allocs := testing.AllocsPerRun(100, func() {
		f.WriteString("I will not allocate when passed a string longer than 32 bytes.\n")
	})
	if allocs != 0 {
		t.Errorf("expected 0 allocs for File.WriteString, got %v", allocs)
	}
}

// Test that it's OK to have parallel I/O and Close on a pipe.
func TestPipeIOCloseRace(t *testing.T) {
	// Skip on wasm, which doesn't have pipes.
	if runtime.GOOS == "js" || runtime.GOOS == "wasip1" || runtime.GOOS == "tamago" {
		t.Skipf("skipping on %s: no pipes", runtime.GOOS)
	}
	t.Parallel()

	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		for {
			n, err := w.Write([]byte("hi"))
			if err != nil {
				// We look at error strings as the
				// expected errors are OS-specific.
				switch {
				case errors.Is(err, ErrClosed),
					strings.Contains(err.Error(), "broken pipe"),
					strings.Contains(err.Error(), "pipe is being closed"),
					strings.Contains(err.Error(), "hungup channel"):
					// Ignore an expected error.
				default:
					// Unexpected error.
					t.Error(err)
				}
				return
			}
			if n != 2 {
				t.Errorf("wrote %d bytes, expected 2", n)
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			var buf [2]byte
			n, err := r.Read(buf[:])
			if err != nil {
				if err != io.EOF && !errors.Is(err, ErrClosed) {
					t.Error(err)
				}
				return
			}
			if n != 2 {
				t.Errorf("read %d bytes, want 2", n)
			}
		}
	}()

	go func() {
		defer wg.Done()

		// Let the other goroutines start. This is just to get
		// a better test, the test will still pass if they
		// don't start.
		time.Sleep(time.Millisecond)

		if err := r.Close(); err != nil {
			t.Error(err)
		}
		if err := w.Close(); err != nil {
			t.Error(err)
		}
	}()

	wg.Wait()
}

// Test that it's OK to call Close concurrently on a pipe.
func TestPipeCloseRace(t *testing.T) {
	// Skip on wasm, which doesn't have pipes.
	if runtime.GOOS == "js" || runtime.GOOS == "wasip1" || runtime.GOOS == "tamago" {
		t.Skipf("skipping on %s: no pipes", runtime.GOOS)
	}
	t.Parallel()

	r, w, err := Pipe()
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	c := make(chan error, 4)
	f := func() {
		defer wg.Done()
		c <- r.Close()
		c <- w.Close()
	}
	wg.Add(2)
	go f()
	go f()
	nils, errs := 0, 0
	for i := 0; i < 4; i++ {
		err := <-c
		if err == nil {
			nils++
		} else {
			errs++
		}
	}
	if nils != 2 || errs != 2 {
		t.Errorf("got nils %d errs %d, want 2 2", nils, errs)
	}
}

func TestRandomLen(t *testing.T) {
	for range 5 {
		dir, err := MkdirTemp(t.TempDir(), "*")
		if err != nil {
			t.Fatal(err)
		}
		base := filepath.Base(dir)
		if len(base) > 10 {
			t.Errorf("MkdirTemp returned len %d: %s", len(base), base)
		}
	}
	for range 5 {
		f, err := CreateTemp(t.TempDir(), "*")
		if err != nil {
			t.Fatal(err)
		}
		base := filepath.Base(f.Name())
		f.Close()
		if len(base) > 10 {
			t.Errorf("CreateTemp returned len %d: %s", len(base), base)
		}
	}
}

func TestCopyFS(t *testing.T) {
	t.Parallel()

	// Test with disk filesystem.
	forceMFTUpdateOnWindows(t, "./testdata/dirfs")
	fsys := DirFS("./testdata/dirfs")
	tmpDir := t.TempDir()
	if err := CopyFS(tmpDir, fsys); err != nil {
		t.Fatal("CopyFS:", err)
	}
	forceMFTUpdateOnWindows(t, tmpDir)
	tmpFsys := DirFS(tmpDir)
	if err := fstest.TestFS(tmpFsys, "a", "b", "dir/x"); err != nil {
		t.Fatal("TestFS:", err)
	}
	if err := verifyCopyFS(t, fsys, tmpFsys); err != nil {
		t.Fatal("comparing two directories:", err)
	}

	// Test whether CopyFS disallows copying for disk filesystem when there is any
	// existing file in the destination directory.
	if err := CopyFS(tmpDir, fsys); !errors.Is(err, fs.ErrExist) {
		t.Errorf("CopyFS should have failed and returned error when there is"+
			"any existing file in the destination directory (in disk filesystem), "+
			"got: %v, expected any error that indicates <file exists>", err)
	}

	// Test with memory filesystem.
	fsys = fstest.MapFS{
		"william":    {Data: []byte("Shakespeare\n")},
		"carl":       {Data: []byte("Gauss\n")},
		"daVinci":    {Data: []byte("Leonardo\n")},
		"einstein":   {Data: []byte("Albert\n")},
		"dir/newton": {Data: []byte("Sir Isaac\n")},
	}
	tmpDir = t.TempDir()
	if err := CopyFS(tmpDir, fsys); err != nil {
		t.Fatal("CopyFS:", err)
	}
	forceMFTUpdateOnWindows(t, tmpDir)
	tmpFsys = DirFS(tmpDir)
	if err := fstest.TestFS(tmpFsys, "william", "carl", "daVinci", "einstein", "dir/newton"); err != nil {
		t.Fatal("TestFS:", err)
	}
	if err := verifyCopyFS(t, fsys, tmpFsys); err != nil {
		t.Fatal("comparing two directories:", err)
	}

	// Test whether CopyFS disallows copying for memory filesystem when there is any
	// existing file in the destination directory.
	if err := CopyFS(tmpDir, fsys); !errors.Is(err, fs.ErrExist) {
		t.Errorf("CopyFS should have failed and returned error when there is"+
			"any existing file in the destination directory (in memory filesystem), "+
			"got: %v, expected any error that indicates <file exists>", err)
	}
}

// verifyCopyFS checks the content and permission of each file inside copied FS to ensure
// the copied files satisfy the convention stipulated in CopyFS.
func verifyCopyFS(t *testing.T, originFS, copiedFS fs.FS) error {
	testDir := filepath.Join(t.TempDir(), "test")
	// umask doesn't apply to the wasip and windows and there is no general way to get masked perm,
	// so create a dir and a file to compare the permission after umask if any
	if err := Mkdir(testDir, ModePerm); err != nil {
		return fmt.Errorf("mkdir %q failed: %v", testDir, err)
	}
	dirStat, err := Stat(testDir)
	if err != nil {
		return fmt.Errorf("stat dir %q failed: %v", testDir, err)
	}
	wantDirMode := dirStat.Mode()

	f, err := Create(filepath.Join(testDir, "tmp"))
	if err != nil {
		return fmt.Errorf("open %q failed: %v", filepath.Join(testDir, "tmp"), err)
	}
	defer f.Close()
	wantFileRWStat, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat file %q failed: %v", f.Name(), err)
	}
	wantFileRWMode := wantFileRWStat.Mode()

	return fs.WalkDir(originFS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			// the dir . is not the dir created by CopyFS so skip checking its permission
			if d.Name() == "." {
				return nil
			}

			dinfo, err := fs.Stat(copiedFS, path)
			if err != nil {
				return err
			}

			if dinfo.Mode() != wantDirMode {
				return fmt.Errorf("dir %q mode is %v, want %v",
					d.Name(), dinfo.Mode(), wantDirMode)
			}
			return nil
		}

		fInfo, err := originFS.Open(path)
		if err != nil {
			return err
		}
		defer fInfo.Close()
		copiedInfo, err := copiedFS.Open(path)
		if err != nil {
			return err
		}
		defer copiedInfo.Close()

		// verify the file contents are the same
		data, err := io.ReadAll(fInfo)
		if err != nil {
			return err
		}
		newData, err := io.ReadAll(copiedInfo)
		if err != nil {
			return err
		}
		if !bytes.Equal(data, newData) {
			return fmt.Errorf("file %q content is %s, want %s", path, newData, data)
		}

		fStat, err := fInfo.Stat()
		if err != nil {
			return err
		}
		copiedStat, err := copiedInfo.Stat()
		if err != nil {
			return err
		}

		// check whether the execute permission is inherited from original FS

		if copiedStat.Mode()&0111&wantFileRWMode != fStat.Mode()&0111&wantFileRWMode {
			return fmt.Errorf("file %q execute mode is %v, want %v",
				path, copiedStat.Mode()&0111, fStat.Mode()&0111)
		}

		rwMode := copiedStat.Mode() &^ 0111 // unset the executable permission from file mode
		if rwMode != wantFileRWMode {
			return fmt.Errorf("file %q rw mode is %v, want %v",
				path, rwMode, wantFileRWStat.Mode())
		}
		return nil
	})
}

func TestCopyFSWithSymlinks(t *testing.T) {
	// Test it with absolute and relative symlinks that point inside and outside the tree.
	testenv.MustHaveSymlink(t)

	// Create a directory and file outside.
	tmpDir := t.TempDir()
	outsideDir := filepath.Join(tmpDir, "copyfs_out")
	if err := Mkdir(outsideDir, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	outsideFile := filepath.Join(outsideDir, "file.out.txt")

	if err := WriteFile(outsideFile, []byte("Testing CopyFS outside"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Create a directory and file inside.
	insideDir := filepath.Join(tmpDir, "copyfs_in")
	if err := Mkdir(insideDir, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	insideFile := filepath.Join(insideDir, "file.in.txt")
	if err := WriteFile(insideFile, []byte("Testing CopyFS inside"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	// Create directories for symlinks.
	linkInDir := filepath.Join(insideDir, "in_symlinks")
	if err := Mkdir(linkInDir, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}
	linkOutDir := filepath.Join(insideDir, "out_symlinks")
	if err := Mkdir(linkOutDir, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	// First, we create the absolute symlink pointing outside.
	outLinkFile := filepath.Join(linkOutDir, "file.abs.out.link")
	if err := Symlink(outsideFile, outLinkFile); err != nil {
		t.Fatalf("Symlink: %v", err)
	}

	// Then, we create the relative symlink pointing outside.
	relOutsideFile, err := filepath.Rel(filepath.Join(linkOutDir, "."), outsideFile)
	if err != nil {
		t.Fatalf("filepath.Rel: %v", err)
	}
	relOutLinkFile := filepath.Join(linkOutDir, "file.rel.out.link")
	if err := Symlink(relOutsideFile, relOutLinkFile); err != nil {
		t.Fatalf("Symlink: %v", err)
	}

	// Last, we create the relative symlink pointing inside.
	relInsideFile, err := filepath.Rel(filepath.Join(linkInDir, "."), insideFile)
	if err != nil {
		t.Fatalf("filepath.Rel: %v", err)
	}
	relInLinkFile := filepath.Join(linkInDir, "file.rel.in.link")
	if err := Symlink(relInsideFile, relInLinkFile); err != nil {
		t.Fatalf("Symlink: %v", err)
	}

	// Copy the directory tree and verify.
	forceMFTUpdateOnWindows(t, insideDir)
	fsys := DirFS(insideDir)
	tmpDupDir := filepath.Join(tmpDir, "copyfs_dup")
	if err := Mkdir(tmpDupDir, 0755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	// TODO(panjf2000): symlinks are currently not supported, and a specific error
	// 			will be returned. Verify that error and skip the subsequent test,
	//			revisit this once #49580 is closed.
	if err := CopyFS(tmpDupDir, fsys); !errors.Is(err, ErrInvalid) {
		t.Fatalf("got %v, want ErrInvalid", err)
	}
	t.Skip("skip the subsequent test and wait for #49580")

	forceMFTUpdateOnWindows(t, tmpDupDir)
	tmpFsys := DirFS(tmpDupDir)
	if err := fstest.TestFS(tmpFsys, "file.in.txt", "out_symlinks/file.abs.out.link", "out_symlinks/file.rel.out.link", "in_symlinks/file.rel.in.link"); err != nil {
		t.Fatal("TestFS:", err)
	}
	if err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		fi, err := d.Info()
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".link" {
			if fi.Mode()&ModeSymlink == 0 {
				return errors.New("original file " + path + " should be a symlink")
			}
			tmpfi, err := fs.Stat(tmpFsys, path)
			if err != nil {
				return err
			}
			if tmpfi.Mode()&ModeSymlink != 0 {
				return errors.New("copied file " + path + " should not be a symlink")
			}
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		newData, err := fs.ReadFile(tmpFsys, path)
		if err != nil {
			return err
		}
		if !bytes.Equal(data, newData) {
			return errors.New("file " + path + " contents differ")
		}

		var target string
		switch fileName := filepath.Base(path); fileName {
		case "file.abs.out.link", "file.rel.out.link":
			target = outsideFile
		case "file.rel.in.link":
			target = insideFile
		}
		if len(target) > 0 {
			targetData, err := ReadFile(target)
			if err != nil {
				return err
			}
			if !bytes.Equal(targetData, newData) {
				return errors.New("file " + path + " contents differ from target")
			}
		}

		return nil
	}); err != nil {
		t.Fatal("comparing two directories:", err)
	}
}

func TestAppendDoesntOverwrite(t *testing.T) {
	testMaybeRooted(t, func(t *testing.T, r *Root) {
		name := "file"
		if err := WriteFile(name, []byte("hello"), 0666); err != nil {
			t.Fatal(err)
		}
		var f *File
		var err error
		if r == nil {
			f, err = OpenFile(name, O_APPEND|O_WRONLY, 0)
		} else {
			f, err = r.OpenFile(name, O_APPEND|O_WRONLY, 0)
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, err := f.Write([]byte(" world")); err != nil {
			f.Close()
			t.Fatal(err)
		}
		if err := f.Close(); err != nil {
			t.Fatal(err)
		}
		got, err := ReadFile(name)
		if err != nil {
			t.Fatal(err)
		}
		want := "hello world"
		if string(got) != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

func TestRemoveReadOnlyFile(t *testing.T) {
	testMaybeRooted(t, func(t *testing.T, r *Root) {
		if err := WriteFile("file", []byte("1"), 0); err != nil {
			t.Fatal(err)
		}
		var err error
		if r == nil {
			err = Remove("file")
		} else {
			err = r.Remove("file")
		}
		if err != nil {
			t.Fatalf("Remove read-only file: %v", err)
		}
		if _, err := Stat("file"); !IsNotExist(err) {
			t.Fatalf("Stat read-only file after removal: %v (want IsNotExist)", err)
		}
	})
}

func TestOpenFileDevNull(t *testing.T) {
	// See https://go.dev/issue/71752.
	t.Parallel()

	f, err := OpenFile(DevNull, O_WRONLY|O_CREATE|O_TRUNC, 0o644)
	if err != nil {
		t.Fatalf("OpenFile(DevNull): %v", err)
	}
	f.Close()
}
