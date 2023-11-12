package umountfs

import (
	"github.com/blang/vfs/mountfs"
	"github.com/blang/vfs"
	"os"
)

type UmountFS struct {
	mfs *mountfs.MountFS
	mounts map[string]vfs.Filesystem
}

type Umounter interface {
	Unmount() error
}

func Create(root vfs.Filesystem) *UmountFS {
	mfs := mountfs.Create(root)
	mounts := make(map[string]vfs.Filesystem, 0)
	u := &UmountFS{mfs, mounts}
	return u
}

func (um UmountFS) Mount(mount vfs.Filesystem, path string) error {
	err := um.mfs.Mount(mount, path)
	if err == nil {
		um.mounts[path] = mount
	}
	return err
}

func (um UmountFS) Unmount() error {
	for _, mount := range um.mounts {
		if m, ok := mount.(Umounter); ok {
			m.Unmount()
		}
	}
	return nil
}

func (um UmountFS) ReadDir(path string) ([]os.FileInfo, error) {
	return um.mfs.ReadDir(path)
}

func (um UmountFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error){
	return um.mfs.OpenFile(name, flag, perm)
}

func (um UmountFS) Mkdir(name string, perm os.FileMode) error {
	return um.mfs.Mkdir(name, perm)
}

func (um UmountFS) Stat(name string) (os.FileInfo, error) {
	return um.mfs.Stat(name)
}


