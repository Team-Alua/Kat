package userfs

import (
	"os"
	"path/filepath"
	"github.com/blang/vfs"
	"github.com/blang/vfs/memfs"
	"github.com/blang/vfs/prefixfs"
	"github.com/Team-Alua/kat/umountfs"
)

func Create(authorId string) (*umountfs.UmountFS, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	gDir := filepath.Join(wd, "global")
	os.Mkdir(gDir, 0777)
	globalFs := prefixfs.Create(vfs.OS(), gDir)
	
	rootFs := memfs.Create()
	rootFs.Mkdir("/tmp", 0777)

	mfs := umountfs.Create(rootFs)
	mfs.Mount(globalFs, "/global")

	// /zips => read only discord uploads 
	ad := filepath.Join(wd, authorId)
	os.Mkdir(ad, 0777)

	authorFs := prefixfs.Create(vfs.OS(), ad)
	mfs.Mount(authorFs, "/local")
	return mfs, nil

}
