package tcpfs

import (
	"github.com/blang/vfs"
	"time"
	"net"
	"os"
)

type TcpFS struct {
	conn net.Conn
}

func Create(source string) (TcpFS, error) {
	conn, err := net.DialTimeout("tcp", source, 1 * time.Second)
	if err != nil {
		return TcpFS{}, err
	}
	return TcpFS{conn}, nil
}

func (f TcpFS) PathSeparator() uint8 {
	return '/';
}


func (f TcpFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	conn := f.conn
	return FileReadWriter{conn}, nil
}

func (f TcpFS) Remove(name string) error {
	return nil

}

func (f TcpFS) Rename(oldpath, newpath string) error {
	return nil
}

func (f TcpFS) Stat(name string) (os.FileInfo, error) {
	return nil, nil
}

func (f TcpFS) Lstat(name string) (os.FileInfo, error) {
	return f.Stat(name)
}

func (f TcpFS) ReadDir(path string) ([]os.FileInfo, error) {
	return nil, nil
}

func (f TcpFS) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (f TcpFS) Unmount() error {
	return f.conn.Close()
}

