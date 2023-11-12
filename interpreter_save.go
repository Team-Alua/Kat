package main

import (
	"archive/zip"
	"github.com/google/uuid"
	"github.com/dop251/goja"
	"strings"
	"path/filepath"
)
type SaveObject struct {}

func (i *Interpreter) LoadSaveBuiltins() {
	vm := i.vm
	saveObj, err := vm.New(vm.Get("Object"))
	if err != nil {
		panic(err)
	}

	vm.Set("save", saveObj)

	saveObj.Set("mount", func(call goja.FunctionCall) goja.Value {
		return i.SaveMount(call)
	})

	saveObj.Set("umount", func(call goja.FunctionCall) goja.Value {
		return i.SaveUmount(call)
	})
}

type MountResult struct {
	BackingZip string
	Path string	
	Err string
}

func (i *Interpreter) SaveMount(call goja.FunctionCall) goja.Value {
	vm := i.vm
	// Generate ftp id
	id := strings.ReplaceAll(uuid.New().String(), "-", "")

	fc := NewFtpClient("10.0.0.5", "2121")
	if resp, ok := fc.Login(); !ok {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}
	defer fc.Kill()

	// This should delete any files
	// on the PS4 on failure
	defer fc.DeleteStage(id)
	// Now upload the files to the PS4
	var archive *zip.ReadCloser
	if resp, ok := fc.UploadSave(archive, id); !ok {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}

	resp, ok := fc.CreateTempFolder(id)

	if !ok {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}
	tmpFolder := resp

	sc := NewSaveClient("10.0.0.5", "1234")
	if resp, ok := sc.Connect(); !ok {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}

	resp, ok = sc.Dump(id, tmpFolder)

	sc.Disconnect()
	if !ok {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}

	targetFolder := filepath.Join("", id)
	// Dump to mount directory
	err := fc.Dump(targetFolder, tmpFolder)
	if err != nil {
		return vm.ToValue(MountResult{Path: "", Err: resp})
	}
	return vm.ToValue(MountResult{BackingZip: "", Path: "/" + id, Err: resp})
}

func (i *Interpreter) SaveUmount(call goja.FunctionCall) goja.Value {
	// Requires mount folder
	// Returns new zip name
	return i.vm.ToValue(nil)
}

