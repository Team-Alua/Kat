package main

import (
	"github.com/dop251/goja"
	"github.com/blang/vfs"
	"github.com/Team-Alua/kat/zipfs"
	"github.com/Team-Alua/kat/ftpfs"
	"os"
)

type MountOptions struct {
	MountType string
	ReadOnly bool
}

func (i *Interpreter) chooseMount(source string, mo MountOptions) (vfs.Filesystem, error){
	mt := mo.MountType
	if mt == "zipfs" {
		// Get file handle
		var flags int
		var size int64
		if mo.ReadOnly {
			flags = os.O_RDONLY
			fi, err := i.fs.Stat(source)
			if err != nil {
				return nil, err
			}
			size = fi.Size()
		} else {
			flags = os.O_CREATE | os.O_EXCL | os.O_WRONLY
		}
		fh, err := i.fs.OpenFile(source, flags,0777)
		if err != nil {
			return nil, err
		}
		return zipfs.Create(fh, size)
	}
	if mt == "ftpfs" {
		return ftpfs.Create("10.0.0.5", "2121", source)
	}
	return nil, nil

}

func (i *Interpreter) Mount(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 3 {
		panic("Must have at least 3 arguments.")
	}

	src, ok1 := fc.Argument(0).Export().(string);
	target, ok2 := fc.Argument(1).Export().(string);
	if !ok1 || !ok2 {
		panic("Invalid file paths.")
	}
	var opts *MountOptions = &MountOptions{}
	err := i.vm.ExportTo(fc.Argument(2), opts)
	if err != nil {
		panic(err)
	}

	if f, err := i.fs.Stat(target); err == nil {
		panic("File might already exists " + f.Name());
	}


	vfs, err := i.chooseMount(src, *opts)
	if err != nil {
		panic(err)
	}
	i.fs.Mkdir(target, 0777)
	if err := i.fs.Mount(vfs, target); err != nil {
		panic(err)
	}
	return i.vm.ToValue(vfs)
}

func (i *Interpreter) Umount(fc goja.FunctionCall) goja.Value {
	return i.vm.ToValue(nil)
}

