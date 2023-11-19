package ftpfs

import (
	"os"
	"github.com/jlaffaye/ftp"
	"github.com/blang/vfs"
	"time"
	"path"
	"io"
	"path/filepath"
)

type FtpFS struct {
	conn *ftp.ServerConn
	remotePath string
}

func Create(ip, port, remotePath string) (FtpFS, error) {
	conn, err := ftp.Dial(ip + ":" + port, ftp.DialWithTimeout(1*time.Second))
	if err != nil {
		return FtpFS{}, err
	}

	err = conn.Login("anonymous", "anonymous")
	if err != nil {
		return FtpFS{}, err
	}
	return FtpFS{conn, remotePath}, nil
}

func (f FtpFS) PathSeparator() uint8 {
	return '/';
}


func (f FtpFS) absPath(relPath string) string {
	return path.Join(f.remotePath, relPath)
}

func (f FtpFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	conn := f.conn
	fp := f.absPath(name)
	if (flag & os.O_WRONLY == os.O_WRONLY) {
		r,w := io.Pipe()
		go func() {
			err := conn.Stor(fp,r)
			if err != nil {
				w.CloseWithError(err)
			} else {
				w.Close()
			}
		}()
		// ftp.conn.Store	
		// Create a pipe
		return FileWriter{File{path: fp, conn: conn}, w}, nil
	} else if (flag & os.O_RDONLY == os.O_RDONLY) {
		r, err := conn.Retr(name)
		if err != nil {
			return nil, err
		}
		// ftp.conn.Retr
		// pass r
		return FileReader{File{path: fp, conn: conn}, r}, nil
	}
	return nil, nil
}

func (f FtpFS) Remove(name string) error {
	fi, err := f.Stat(name);
	if err != nil {
		return err
	}

	if fi.IsDir() {
		f.conn.ChangeDir(filepath.Dir(name))
		f.conn.RemoveDir(filepath.Base(name))
	} else {
		f.conn.Delete(name)
	}
	return nil

}

func (f FtpFS) Rename(oldpath, newpath string) error {
	oldpath = f.absPath(oldpath)
	newpath = f.absPath(newpath)
	return f.conn.Rename(oldpath, newpath)
}

func (f FtpFS) Stat(name string) (os.FileInfo, error) {
	name = f.absPath(name)
	entry, err := f.conn.GetEntry(name)
	if err != nil {
		return nil, err
	}
	return fileInfo{entry}, nil
}

func (f FtpFS) Lstat(name string) (os.FileInfo, error) {
	return f.Stat(name)
}

func (f FtpFS) ReadDir(path string) ([]os.FileInfo, error) {
	fp := f.absPath(path)
	entries, err := f.conn.List(fp)
	if err != nil {
		return nil, err
	}
	ret := make([]os.FileInfo, 0)
	for _, entry := range entries {
		if entry.Name == "." || entry.Name == ".." {
			continue
		}
		ret = append(ret, fileInfo{entry})
	}
	return ret, nil
}

func (f FtpFS) Mkdir(name string, perm os.FileMode) error {
	fp := f.absPath(name)
	return f.conn.MakeDir(fp)
}


func (f FtpFS) Unmount() error {
	return f.conn.Quit()
}

