package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"./linkfs"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Printf("usage: %s MOUNTPOINT ORIGINAL\n", path.Base(os.Args[0]))
		os.Exit(2)
	}
	orig := flag.Arg(1)
	linkfilesystem := linkfs.NewLinkFs(orig)
	nfs := pathfs.NewPathNodeFs(linkfilesystem, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(0), nfs.Root(), nil)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}
	server.Serve()
}
