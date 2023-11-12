package main

import (
	"github.com/dop251/goja"
	"github.com/Team-Alua/kat/umountfs"
	"fmt"
)


type Interpreter struct {
	vm *goja.Runtime
	rw *DiscordReadWriter
	fs *umountfs.UmountFS
}

func NewInterpreter(rw *DiscordReadWriter, fs *umountfs.UmountFS) *Interpreter {
	i := &Interpreter{}	
	i.vm = goja.New()
	i.rw = rw
	i.fs = fs
	return i
}

func (i *Interpreter) Run(name, code string) error {
	vm := i.vm
	i.LoadBuiltins()
	defer func() {
//		if err := recover(); err != nil {
//			fmt.Println("Error", err);
//		}
		i.fs.Unmount()
		// Cleanup fs
	}()
	_, err := vm.RunScript(name, code)
	return err
}

func (i *Interpreter) LoadBuiltins() {
	i.LoadFsBuiltins()
	i.LoadDiscordBuiltins()
	i.LoadSaveBuiltins()
	i.LoadHttpBuiltins()
	vm := i.vm


	vm.Set("log", func(call goja.FunctionCall) goja.Value {
		for _, arg := range call.Arguments {
			fmt.Printf("%s", arg)
		}
		fmt.Printf("\n")
		return i.vm.ToValue(nil)
	});

	vm.Set("run", func(script string) goja.Value {
		vm.Interrupt("run " + script);
		return vm.ToValue(nil)
	});

	vm.Set("exit", func(call goja.FunctionCall) goja.Value {
		vm.Interrupt("exit")
		return vm.ToValue(nil)
	});
}

