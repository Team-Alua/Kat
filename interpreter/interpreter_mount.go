package interpreter

import (
    "github.com/dop251/goja"
    "github.com/blang/vfs"
    "github.com/Team-Alua/kat/zipfs"
    "github.com/Team-Alua/kat/tcpfs"
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
    if mt == "tcpfs" {
        return tcpfs.Create()
    }
    return nil, nil

}

func (i *Interpreter) Mount(fc goja.FunctionCall) goja.Value {
    vm := i.vm
    if len(fc.Arguments) < 3 {
        vm.Interrupt("Must have at least 3 arguments.")
        return vm.ToValue(nil)
    }

    src, ok1 := fc.Argument(0).Export().(string);
    target, ok2 := fc.Argument(1).Export().(string);
    if !ok1 || !ok2 {
        vm.Interrupt("Invalid file paths.")
        return vm.ToValue(nil)
    }
    var opts *MountOptions = &MountOptions{}
    err := i.vm.ExportTo(fc.Argument(2), opts)
    if err != nil {
        vm.Interrupt(err)
        return vm.ToValue(nil)
    }

    if f, err := i.fs.Stat(target); err == nil {
        vm.Interrupt("File might already exists " + f.Name());
        return vm.ToValue(nil)
    }


    vfs, err := i.chooseMount(src, *opts)
    if err != nil {
        vm.Interrupt(err)
        return vm.ToValue(nil)
    }

    if err := i.fs.Mount(vfs, target); err != nil {
        vm.Interrupt(err)
        return vm.ToValue(nil)
    }
    return i.vm.ToValue(vfs)
}

func (i *Interpreter) Umount(fc goja.FunctionCall) goja.Value {
    if len(fc.Arguments) < 1 {
        panic("Must have at least 1 argument.")
    }
    var p string
    if err := i.vm.ExportTo(fc.Argument(0), &p); err != nil {
        panic(err)
    }

    if err := i.fs.Unmount(p); err != nil {
        panic(err)
    }

    return i.vm.ToValue(nil)
}

