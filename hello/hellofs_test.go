package hellofs

import (
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

var h HelloFs

func TestGetAttr(t *testing.T) {
	type test struct {
		name   string
		attr   fuse.Attr
		status fuse.Status
	}

	tests := []test{
		//{"doesnotexist", nil, fuse.ENOENT},
		{"file.txt", fuse.Attr{Mode: fuse.S_IFREG | 0644, Size: uint64(8)}, fuse.OK},
		{"", fuse.Attr{Mode: fuse.S_IFDIR | 0755, Size: 0}, fuse.OK},
	}

	for _, testcase := range tests {

		attribute, status := h.GetAttr(testcase.name, nil)
		if *attribute != testcase.attr {
			t.Errorf("Attributes did not match, got: %v, expected: %v", *attribute, testcase.attr)
		}
		if status != testcase.status {
			t.Errorf("Status code did not match, got: %d, expected: %d", status, testcase.status)
		}
	}
}

func TestOpenDir(t *testing.T) {
	type test struct {
		name    string
		entries []fuse.DirEntry
		status  fuse.Status
	}

	tests := []test{
		{"", []fuse.DirEntry{
			{Name: "file.txt", Mode: fuse.S_IFREG},
			{Name: "folder", Mode: fuse.S_IFDIR},
		}, fuse.OK},
		{"folder", []fuse.DirEntry{{Name: "inner_file.txt", Mode: fuse.S_IFREG}}, fuse.OK},
		{"nonexistent", nil, fuse.ENOENT},
	}

	for _, testcase := range tests {
		c, status := h.OpenDir(testcase.name, nil)
		if testcase.status != status {
			t.Errorf("Status code did not match, got: %d, expected: %d", status, testcase.status)
		}

		entriesMatch := true
		for i, entry := range c {
			if entry.Name != testcase.entries[i].Name || entry.Mode != testcase.entries[i].Mode {
				entriesMatch = false
			}
		}
		if !entriesMatch {
			t.Errorf("Entries did not match, got: %v, expected: %v", c, testcase.entries)
		}
	}
}

func TestOpen(t *testing.T) {
	type test struct {
		name   string
		file   nodefs.File
		status fuse.Status
	}
	tests := []test{
		{"file.txt", nil, fuse.OK},
		{"nonexistent", nil, fuse.ENOENT},
	}

	for _, testcase := range tests {
		_, status := h.Open(testcase.name, 0, nil)
		if status != testcase.status {
			t.Errorf("Return Status did not match, got: %d, expected %d", status, testcase.status)
		}
	}
}
