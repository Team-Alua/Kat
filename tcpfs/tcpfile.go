package tcpfs

import (
	"os"
	"net"
	"time"
	filesystem "io/fs"
)
type fileInfo struct {
	name string
	dir bool
	size int64
}

func (f fileInfo) ModTime() time.Time {
	return time.Time{}
}

func (f fileInfo) Name() string {
	return f.name
}

func (f fileInfo) Size() int64 {
	return f.size
}

func (f fileInfo) Mode() filesystem.FileMode {
	return 0
}

func (f fileInfo) IsDir() bool {
	return f.dir
}

func (f fileInfo) Sys() any {
	return nil
}

type FileReadWriter struct {
	conn net.Conn
}

func (f FileReadWriter) Name() string {
	return ""
}

func (f FileReadWriter) Stat() (os.FileInfo, error) {
	return nil, nil
}

func (f FileReadWriter) Read(p []byte) (int, error) {
	return f.conn.Read(p)
}

func (f FileReadWriter) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

func (f FileReadWriter) Write(p []byte) (int, error) {
	return f.conn.Write(p)

}

func (f FileReadWriter) Truncate(size int64) (err error) {
	return nil
}

func (f FileReadWriter) Seek(offset int64, whence int) (n int64, err error) {
	return 0, nil
}

func (f FileReadWriter) Sync() error {
	return nil
}


func (f FileReadWriter) Close() error {
	return f.conn.Close()
}

