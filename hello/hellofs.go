package hellofs

import (
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type HelloFs struct {
	pathfs.FileSystem
}

func (me *HelloFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	switch name {
	case "file.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0755,
		}, fuse.OK
	case "folder":
		return &fuse.Attr{
			Mode: fuse.S_IFDIR | 0644, Size: uint64(len(name)),
		}, fuse.OK
	case "folder/inner_file.txt":
		return &fuse.Attr{
			Mode: fuse.S_IFREG | 0644, Size: uint64(len(name)),
		}, fuse.OK
	}
	return nil, fuse.ENOENT
}

func (me *HelloFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	if name == "" {
		c = []fuse.DirEntry{
			{Name: "file.txt", Mode: fuse.S_IFREG},
			{Name: "folder", Mode: fuse.S_IFDIR},
		}
		return c, fuse.OK
	}
	if name == "folder" {
		c = []fuse.DirEntry{
			{Name: "inner_file.txt", Mode: fuse.S_IFREG},
		}
		return c, fuse.OK
	}

	return nil, fuse.ENOENT
}

func (me *HelloFs) Open(name string, flags uint32, context *fuse.Context) (file nodefs.File, code fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nil, fuse.EPERM
	}
	if name == "file.txt" {
		return nodefs.NewDataFile([]byte("Hello I am the outer file!")), fuse.OK
	}
	if name == "folder/inner_file.txt" {
		return nodefs.NewDataFile([]byte("Hello, I am the inner file!")), fuse.OK
	}
	return nil, fuse.ENOENT
}
