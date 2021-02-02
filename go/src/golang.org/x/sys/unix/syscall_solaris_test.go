// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build solaris

package unix_test

import (
	"os/exec"
	"testing"

	"golang.org/x/sys/unix"
)

func TestPipe2(t *testing.T) {
	const s = "hello"
	var pipes [2]int
	err := unix.Pipe2(pipes[:], 0)
	if err != nil {
		t.Fatalf("pipe2: %v", err)
	}
	r := pipes[0]
	w := pipes[1]
	go func() {
		n, err := unix.Write(w, []byte(s))
		if err != nil {
			t.Errorf("bad write: %v", err)
			return
		}
		if n != len(s) {
			t.Errorf("bad write count: %d", n)
			return
		}
		err = unix.Close(w)
		if err != nil {
			t.Errorf("bad close: %v", err)
			return
		}
	}()
	var buf [10 + len(s)]byte
	n, err := unix.Read(r, buf[:])
	if err != nil {
		t.Fatalf("bad read: %v", err)
	}
	if n != len(s) {
		t.Fatalf("bad read count: %d", n)
	}
	if string(buf[:n]) != s {
		t.Fatalf("bad contents: %s", string(buf[:n]))
	}
	err = unix.Close(r)
	if err != nil {
		t.Fatalf("bad close: %v", err)
	}
}

func TestStatvfs(t *testing.T) {
	if err := unix.Statvfs("", nil); err == nil {
		t.Fatal(`Statvfs("") expected failure`)
	}

	statvfs := unix.Statvfs_t{}
	if err := unix.Statvfs("/", &statvfs); err != nil {
		t.Errorf(`Statvfs("/") failed: %v`, err)
	}

	if t.Failed() {
		mount, err := exec.Command("mount").CombinedOutput()
		if err != nil {
			t.Logf("mount: %v\n%s", err, mount)
		} else {
			t.Logf("mount: %s", mount)
		}
	}
}
