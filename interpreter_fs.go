package main

import (
	"github.com/dop251/goja"
	"github.com/blang/vfs"
	"os"
	"io/fs"
	"io"
	"bufio"
)

type JSFileInfo struct {
	Name string
	Size int64
	Dir bool
}

func (i *Interpreter) LoadFsIntoInstance(f *goja.Object) {
	vm := i.vm
	f.Set("open", func(fc goja.FunctionCall) goja.Value {
		return i.OpenFile(fc)
	})

	f.Set("mkdir", func(fc goja.FunctionCall) goja.Value {
		return i.Mkdir(fc)
	})

	f.Set("readdir", func(fc goja.FunctionCall) goja.Value {
		files, err := i.fs.ReadDir(fc.Argument(0).Export().(string))
		if err != nil {
			panic(err)
		}
		out := make([]JSFileInfo, 0)
		for _, f := range files {
			of := JSFileInfo{Name: f.Name(), Size: f.Size(), Dir: f.IsDir()}
			out = append(out, of)
		}
		return vm.ToValue(out)
	})

	f.Set("mount", func(fc goja.FunctionCall) goja.Value {
		return i.Mount(fc)
	});

	f.Set("umount", func(fc goja.FunctionCall) goja.Value {
		return i.Umount(fc)
	});

	f.Set("read", func(fc goja.FunctionCall) goja.Value {
		return i.ReadFile(fc)
	})

	f.Set("write", func(fc goja.FunctionCall) goja.Value {
		return i.WriteFile(fc)
	})

	f.Set("writeString", func(fc goja.FunctionCall) goja.Value {
		return i.WriteFileString(fc)
	})

	f.Set("copy", func(fc goja.FunctionCall) goja.Value {
		return i.CopyFile(fc)
	})

	f.Set("close", func(fc goja.FunctionCall) goja.Value {
		return i.CloseFile(fc)
	})

	f.Set("readline", func(fc goja.FunctionCall) goja.Value {
		return i.ReadFileLine(fc)
	})

}

func (i *Interpreter) LoadFsConstants(f *goja.Object) {
	// Exclusive
	f.Set("O_RDONLY", os.O_RDONLY)
	f.Set("O_WRONLY", os.O_WRONLY)
	f.Set("O_RDWR", os.O_RDWR)

	f.Set("O_APPEND", os.O_APPEND)
	f.Set("O_CREATE", os.O_CREATE)
	f.Set("O_EXCL", os.O_EXCL)
	f.Set("O_SYNC", os.O_SYNC)
	f.Set("O_TRUNC", os.O_TRUNC)
}

func (i *Interpreter) LoadFsBuiltins() {
	vm := i.vm
	_ = i.fs
	fsObj, err := vm.New(vm.Get("Object"))
	if err != nil {
		panic(err)
	}
	vm.Set("fs", fsObj)
	i.LoadFsIntoInstance(fsObj)

	fsCstObj, err := vm.New(vm.Get("Object"))
	fsObj.Set("constants", fsCstObj)
	i.LoadFsConstants(fsCstObj)
}

func (i *Interpreter) OpenFile(fc goja.FunctionCall) goja.Value {
	// OpenFile(path, flags, permissions)
	fp, ok := fc.Argument(0).Export().(string);
	if !ok {
		panic("First argument must be a string.")
	}

	flags := int(fc.Argument(1).ToInteger());
	perm := fs.FileMode(fc.Argument(2).ToInteger());
	fh, err := i.fs.OpenFile(fp, flags, perm)
	if err != nil {
		panic(err)
	}
	return i.vm.ToValue(fh)
}

// fs.read(fh, byteCount) => ArrayBuffer
func (i *Interpreter) ReadFile(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 2 {
		panic("Invalid argument count.")
	}

	var r io.Reader
	vm := i.vm
	if err := vm.ExportTo(fc.Argument(0), &r); err != nil {
		panic(err)
	}
	n := fc.Argument(1).ToInteger()
	if n == 0 {
		panic("Invalid read amount.")
	}

	p := make([]byte, n)
	if _, err := r.Read(p); err != nil {
		p = nil
		panic(err)
	}
	return i.vm.ToValue(i.vm.NewArrayBuffer(p))
}


func (i *Interpreter) WriteFile(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 2 {
		panic("Invalid argument count.")
	}

	var w io.Writer
	vm := i.vm
	if err := vm.ExportTo(fc.Argument(0), &w); err != nil {
		panic(err)
	}

	var p []byte
	if err := vm.ExportTo(fc.Argument(1), &p); err != nil {
		panic(err)
	}
	n, err := w.Write(p)
	if err != nil {
		panic(err)
	}
	return i.vm.ToValue(n)
}
func (i *Interpreter) WriteFileString(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 2 {
		panic("Invalid argument count.")
	}

	var w io.Writer
	vm := i.vm
	if err := vm.ExportTo(fc.Argument(0), &w); err != nil {
		panic(err)
	}

	var p string
	if err := vm.ExportTo(fc.Argument(1), &p); err != nil {
		panic(err)
	}
	n, err := w.Write([]byte(p))
	if err != nil {
		panic(err)
	}
	return i.vm.ToValue(n)
}

func (i *Interpreter) CloseFile(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 1 {
		panic("Invalid argument count.")
	}
	var f vfs.File
	vm := i.vm
	if err := vm.ExportTo(fc.Argument(0), &f); err != nil {
		panic(err)
	}

	if err := f.Close(); err != nil {
		panic(err)
	}
	return i.vm.ToValue(nil)
}

func (i *Interpreter) CopyFile(fc goja.FunctionCall) goja.Value {

	if len(fc.Arguments) < 2 {
		panic("Invalid argument count.")
	}
	vm := i.vm

	var w io.Writer
	if err := vm.ExportTo(fc.Argument(0), &w); err != nil {
		panic(err)
	}
	
	var r io.Reader
	if err := vm.ExportTo(fc.Argument(1), &r); err != nil {
		panic(err)
	}

	_, err := io.Copy(w,r)
	if err != nil {
		panic(err)
	}	 
	return i.vm.ToValue(nil)
}

func (i *Interpreter) Mkdir(fc goja.FunctionCall) goja.Value {
	fp, ok := fc.Argument(0).Export().(string);
	if !ok {
		panic("First argument must be a string.")
	}
	var fm os.FileMode
	fm, ok = fc.Argument(1).Export().(os.FileMode);
	if !ok {
		fm = 0777
	}
	if err := i.fs.Mkdir(fp, fm); err != nil {
		panic(err)
	}
	return i.vm.ToValue(nil)
}

func (i *Interpreter) ReadFileLine(fc goja.FunctionCall) goja.Value {
	if len(fc.Arguments) < 1 {
		panic("Invalid argument count.")
	}

	var f vfs.File
	vm := i.vm
	if err := vm.ExportTo(fc.Argument(0), &f); err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	b := scanner.Bytes()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return i.vm.ToValue(string(b))
}
