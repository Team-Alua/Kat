package zipfs

import (
	"archive/zip"
	"path"
	"os"
	"github.com/blang/vfs"
	"errors"
	filesystem "io/fs"
	"strings"
)
var (
	// ErrReadOnly is returned if the file is read-only and write operations are disabled.
	ErrReadOnly = errors.New("File is read-only")
	// ErrWriteOnly is returned if the file is write-only and read operations are disabled.
	ErrWriteOnly = errors.New("File is write-only")
	// ErrIsDirectory is returned if the file under operation is not a regular file but a directory.
	ErrIsDirectory = errors.New("Is directory")

	ErrNotPerm = errors.New("Operation not permitted")

	ErrNotExist = errors.New("File does not exist")
)

func findFile(archive *zip.Reader, name string) (*zip.File, error) {
	for _, file := range archive.File {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, ErrNotExist
}


type ZipFS struct {
	f vfs.File
	r *zip.Reader
	w *zip.Writer
	fm map[string]*fileInfo
}

func Create(f vfs.File, size int64) (ZipFS, error) {
	var r *zip.Reader
	var w *zip.Writer
	fm := make(map[string]*fileInfo)
	if (size == 0) {
		w = zip.NewWriter(f)
	} else {
		nr, err := zip.NewReader(f, size)
		if err != nil {
			return ZipFS{}, err
		}
		r = nr
	}
	fs := ZipFS{f, r, w, fm}
	if (size != 0) {
		fs.buildTree()
	}
	return fs, nil
}

func (f ZipFS) buildTree() {
	for k := range f.fm {
	    delete(f.fm, k)
	}
	// PS4/SAVEDATA
	// gets turned into
	// /PS4
	// /PS4/SAVEDATA

	for _, file := range f.r.File {
		root := "/"
		for _, path := range strings.Split(file.Name, "/") {
			root += path
			if root[1:] == file.Name {
				break
			}
			f.fm[root] = &fileInfo{name: root, dir: true, size: 0}
			root += "/"
		}
		name := file.Name
		if file.FileInfo().IsDir() {
			name = name[0:len(file.Name) - 1]
		}
		fi := &fileInfo{name: "/" + name, dir: file.FileInfo().IsDir()}
		fi.size = file.FileInfo().Size()
		f.fm["/" + name] = fi
		
	}
}

func (f ZipFS) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (f ZipFS) PathSeparator() uint8 {
	return '/';
}

func (f ZipFS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	// Remove leading /
	name = name[1:]
	if (flag & os.O_WRONLY == os.O_WRONLY) {
		// Use writer
		w, err := f.w.Create(name)
		if err != nil {
			return nil, err
		}
		return File{w: w}, nil
	} else if (flag & os.O_RDONLY == os.O_RDONLY) {
		// Use reader
		rc, err := f.r.Open(name)
		if err != nil {
			return nil, err
		}
		return File{rc: rc}, nil
	}

	// Not supported
	return nil, ErrNotPerm
}

func (f ZipFS) Remove(name string) error {
	return ErrNotPerm
}

func (f ZipFS) Rename(oldpath, newpath string) error {
	return ErrNotPerm

}

func (f ZipFS) Stat(name string) (filesystem.FileInfo, error) {
	name = path.Clean(name)
	file, err := findFile(f.r, name)
	if err != nil {
		return nil, err
	}
	return file.FileInfo(), nil
}

func (f ZipFS) Lstat(name string) (os.FileInfo, error) {
	return f.Stat(name)
}

func (f ZipFS) ReadDir(fp string) (entries []filesystem.FileInfo, err error) {
	entries = make([]filesystem.FileInfo, 0)
	for _, file := range f.fm {
		if path.Dir(file.Name()) != fp {
			continue
		}
		// Return only the names
		name := path.Base(file.Name())
		dir := file.IsDir()
		size := file.Size()
		entries = append(entries, fileInfo{name, dir, size})
	}
	return 
}

func (f ZipFS) Unmount() error {
	if f.w != nil {
		f.w.Close()
	}
	return f.f.Close()
}
