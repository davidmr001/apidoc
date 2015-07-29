// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package scanner

import (
	"testing"

	"github.com/issue9/assert"
)

func TestScanner_next(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("ab\ncd"),
	}

	a.Equal('a', s.next())
	a.Equal('b', s.next())
	a.Equal('\n', s.next())
	a.Equal('c', s.next())
	a.Equal('d', s.next())
	a.Equal(eof, s.next())
	a.Equal(eof, s.next())
}

func TestScanner_match(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("ab\ncd"),
	}

	a.False(s.match("b")).Equal(0, s.pos)
	a.True(s.match("a")).Equal(1, s.pos)

	s.backup()
	s.backup()
	a.True(s.match("a")).Equal(1, s.pos)
	a.True(s.match("b")).Equal(2, s.pos)

	s.pos = len(s.data)
	a.False(s.match("ab"))
}

func TestScanner_lineNumber(t *testing.T) {
	a := assert.New(t)

	s := &scanner{
		data: []byte("adf\n\nadf"),
		pos:  3,
	}
	a.Equal(0, s.lineNumber())

	s.pos = 4
	a.Equal(1, s.lineNumber())
}

func TestScanner_scan(t *testing.T) {
	a := assert.New(t)
	s, err := newScanner(cstyle)
	a.NotError(err).NotNil(s)

	a.NotError(s.scan("./testcode/php1.php"))

	php1, found := s.tree.Docs["php1"]
	a.True(found).NotNil(php1)

	a.Equal(php1[0].Methods, "get").
		Equal(php1[0].URL, "/api/php1/get")
}

func TestScan(t *testing.T) {
	a := assert.New(t)

	tree, err := Scan("./testcode", true, "", nil)
	a.NotError(err).NotNil(tree)

	php1, found := tree.Docs["php1"]
	a.True(found).NotNil(php1)

	a.Equal(php1[0].Methods, "get").
		Equal(php1[0].URL, "/api/php1/get")

	php2, found := tree.Docs["php2"]
	a.True(found).NotNil(php2)

	for _, v := range php2 {
		println(v.URL)
	}
	a.Equal(php2[0].Methods, "get").
		Equal(php2[0].URL, "/api/php2/get")
}

func TestDetectLangType(t *testing.T) {
	a := assert.New(t)

	l, err := detectLangType([]string{".abc1", ".abc1", ".abc1"})
	a.Error(err).Equal(0, len(l))

	l, err = detectLangType([]string{".js", ".php", ".abc1"})
	a.NotError(err).Equal("js", l)
}

func TestDetectDirLangType(t *testing.T) {
	a := assert.New(t)

	l, err := detectDirLangType("./")
	a.NotError(err).Equal(l, "go")

	l, err = detectDirLangType("./testdir")
	a.Error(err).Equal(0, len(l))
}

func TestRecursivePath(t *testing.T) {
	a := assert.New(t)

	paths, err := recursivePath("./testdir", false, ".1", ".2")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	paths, err = recursivePath("./testdir", true, ".1", ".2")
	a.NotError(err)
	a.Contains(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
		"testdir/testfile.2",
	})

	paths, err = recursivePath("./testdir/testdir1", true, ".1", ".2")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir1/testfile.2",
	})

	paths, err = recursivePath("./testdir", true, ".1")
	a.NotError(err)
	a.Equal(paths, []string{
		"testdir/testdir1/testfile.1",
		"testdir/testdir2/testfile.1",
		"testdir/testfile.1",
	})
}