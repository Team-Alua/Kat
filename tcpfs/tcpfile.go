package tcpfs

import (
	"os"
	"net"
)

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
	f.conn = nil
	return nil
}

