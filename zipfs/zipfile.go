package zipfs
import (
    filesystem "io/fs"
    "time"
    "io"
    "os"
    "errors"
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

type File struct {
    rc io.ReadCloser
    w io.Writer
}

func (f File) Name() string {
    return ""
}

func (f File) Stat() (os.FileInfo, error) {
    return fileInfo{}, nil
}

func (f File) Read(p []byte) (int, error) {
    if f.rc == nil {
        return 0, errors.New("Invalid operation")
    }
    return f.rc.Read(p)
}

func (f File) ReadAt(p []byte, off int64) (n int, err error) {
    return 0, nil
}

func (f File) Write(p []byte) (int, error) {
    if f.w == nil {
        return 0, errors.New("Invalid operation")
    }
    return f.w.Write(p)

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
    if f.rc != nil {
        return f.rc.Close()
    }
    return nil
}

