package ftpfs

import (
	"os"
	"github.com/jlaffaye/ftp"
	"time"
	"io"
	filesystem "io/fs"
	"errors"
)

type fileInfo struct {
	entry *ftp.Entry
}

func (f fileInfo) ModTime() time.Time {
	return f.entry.Time
}

func (f fileInfo) Name() string {
	return f.entry.Name
}

func (f fileInfo) Size() int64 {
	return int64(f.entry.Size)
}

func (f fileInfo) Mode() filesystem.FileMode {
	return 0
}

func (f fileInfo) IsDir() bool {
	return f.entry.Type == ftp.EntryTypeFolder
}

func (f fileInfo) Sys() any {
	return nil
}

type File struct {
	conn *ftp.ServerConn
	path string
}

func (f File) Name() string {
	return ""
}

func (f File) Stat() (os.FileInfo, error) {
	return fileInfo{}, nil
}

func (f File) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}


func (f File) Truncate(size int64) (err error) {
	return nil
}

func (f File) Seek(offset int64, whence int) (n int64, err error) {
	return 0, nil
}

func (f File) Sync() error {
	return nil
}

func (f File) Close() error {
	return nil
}

type FileReader struct {
	File	
	r io.ReadCloser
}

func (f FileReader) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

func (f FileReader) Write(p []byte) (int, error) {
	return 0, errors.New("Invalid operation")
}

func (f FileReader) Close() error {
	return f.r.Close()
}

type FileWriter struct {
	File
	w io.Writer
}

func (f FileWriter) Read(p []byte) (int, error) {
	return 0, errors.New("Invalid operation")
}


func (f FileWriter) Write(p []byte) (int, error) {
	return f.w.Write(p)
}
