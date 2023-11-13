package tcpfs

import (
	"github.com/blang/vfs"
	"time"
	"net"
	"os"
)

// TODO: Automatically clean up file handles on close
type TcpFS struct {
}

func Create() (TcpFS, error) {
	return TcpFS{}, nil
}

func (f TcpFS) PathSeparator() uint8 {
	return '/';
}


func (f TcpFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	conn, err := net.DialTimeout("tcp", name[1:], 1 * time.Second)
	if err != nil {
		return FileReadWriter{}, err
	}
	return FileReadWriter{conn}, nil
}

func (f TcpFS) Remove(name string) error {
	return nil

}

func (f TcpFS) Rename(oldpath, newpath string) error {
	return nil
}

func (f TcpFS) Stat(name string) (os.FileInfo, error) {
	return &fileInfo{name: name, dir: false, size: 0}, nil
}

func (f TcpFS) Lstat(name string) (os.FileInfo, error) {
	return f.Stat(name)
}

func (f TcpFS) ReadDir(path string) ([]os.FileInfo, error) {
	return make([]os.FileInfo, 0), nil
}

func (f TcpFS) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (f TcpFS) Unmount() error {
	return nil
}

