package linkfs

import (
	"os"
	"syscall"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
)

var cwd, _ = os.Getwd()

var l = NewLinkFs(cwd + "/orig_test")

func TestGetAttr(t *testing.T) {

	tests := []string{
		"file.py",
		"folder",
		"folder/file",
	}

	for _, testcase := range tests {
		st := syscall.Stat_t{}
		syscall.Lstat(cwd+"/orig_test/"+testcase, &st)
		a := &fuse.Attr{}
		a.FromStat(&st)

		attribute, status := l.GetAttr(testcase, nil)

		if *attribute != *a {
			t.Errorf("Attributes did not match, got: %v, expected: %v", *attribute, *a)
		}

		if status != fuse.OK {
			t.Errorf("Status code did not match, got: %d, expected: %d", status, fuse.OK)
		}
	}

}

func TestOpenDir(t *testing.T) {

	type test struct {
		name   string
		status fuse.Status
	}
	tests := []test{
		{"", fuse.OK},
		{"folder", fuse.OK},
		{"nonexistent", fuse.ENOENT},
	}

	for _, testcase := range tests {
		_, status := l.OpenDir(testcase.name, nil)

		if testcase.status != status {
			t.Errorf("Status code did not match, got: %d, expected: %d", status, testcase.status)
		}
	}
}

func TestOpen(t *testing.T) {
	type test struct {
		name   string
		status fuse.Status
	}

	tests := []test{
		{"file.py", fuse.OK},
		{"folder/file", fuse.OK},
		{"doesnotexist.txt", fuse.ENOENT},
	}

	for _, testcase := range tests {
		file, status := l.Open(testcase.name, 0, nil)

		if testcase.status != status {
			t.Errorf("Status code did not match, got: %d, expected: %d", status, testcase.status)
		}

		if status == fuse.OK {
			expected_file := "readOnlyFile(loopbackFile(" + cwd + "/orig_test/" + testcase.name + "))"
			if expected_file != file.String() {
				t.Errorf("Files did not match, expected: %s, got %s", expected_file, file.String())
			}
		}
	}
}
