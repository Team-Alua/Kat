package main

import (
	"os"
	"github.com/dop251/goja"
)

func (i *Interpreter) LoadHttpBuiltins() {
	vm := i.vm
	vm.Set("download", func(path, url string) goja.Value {
		return i.DownloadFile(path, url)
	});
}


func (i *Interpreter) DownloadFile(path, url string) goja.Value {
	fs := i.fs
	fh, err := fs.OpenFile(path, os.O_CREATE | os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer fh.Close()
	DownloadToWriter(fh, url)
	return i.vm.ToValue(nil)
}
