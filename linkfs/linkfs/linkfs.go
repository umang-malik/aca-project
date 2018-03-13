package linkfs

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"syscall"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

type LinkFs struct {
	pathfs.FileSystem
	Root string
}

func NewLinkFs(root string) pathfs.FileSystem {
	root, err := filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	return &LinkFs{
		FileSystem: pathfs.NewDefaultFileSystem(),
		Root:       root,
	}
}

func (fs *LinkFs) Write(name string, context *fuse.Context) fuse.Status {
	return fuse.ENOENT
}
func (fs *LinkFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	fullpath := fs.GetPath(name)
	fmt.Println(fullpath)
	st := syscall.Stat_t{}
	var err error = nil
	if name == "" {
		err = syscall.Stat(fullpath, &st)
	} else {
		err = syscall.Lstat(fullpath, &st)
	}

	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	a := &fuse.Attr{}
	a.FromStat(&st)
	return a, fuse.OK
}

func (fs *LinkFs) OpenDir(name string, context *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	f, err := os.Open(fs.GetPath(name))
	defer f.Close()
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	capacity := 499
	output := make([]fuse.DirEntry, 0, capacity)
	for {
		infos, err := f.Readdir(capacity)
		for i := range infos {
			//Handle https://code.google.com/p/go/issues/detail?id=5960
			if infos[i] == nil {
				continue
			}
			n := infos[i].Name()
			d := fuse.DirEntry{
				Name: n,
			}
			if s := fuse.ToStatT(infos[i]); s != nil {
				d.Mode = uint32(s.Mode)
				d.Ino = s.Ino
			} else {
				log.Printf("ReadDir entry %q for %q has no stat info", n, name)
			}
			output = append(output, d)
		}
		if len(infos) < capacity || err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Readdir() returned err:", err)
			break
		}
	}
	fmt.Println(output)
	return output, fuse.OK
}

func (fs *LinkFs) Open(name string, flags uint32, context *fuse.Context) (fuseFile nodefs.File, status fuse.Status) {
	f, err := os.OpenFile(fs.GetPath(name), int(flags), 0)
	if err != nil {
		return nil, fuse.ToStatus(err)
	}
	osFile := nodefs.NewLoopbackFile(f)
	fmt.Println(nodefs.NewReadOnlyFile(osFile))
	return nodefs.NewReadOnlyFile(osFile), fuse.OK
}

func (fs *LinkFs) GetPath(relPath string) string {
	return filepath.Join(fs.Root, relPath)
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("usage: %s MOUNTPOINT ORIGINAL\n", path.Base(os.Args[0]))
		os.Exit(2)
	}
	orig := flag.Arg(1)
	linkfs := NewLinkFs(orig)
	nfs := pathfs.NewPathNodeFs(linkfs, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
