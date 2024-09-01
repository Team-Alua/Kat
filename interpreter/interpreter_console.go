package interpreter

import (
    "github.com/dop251/goja"
    "fmt"
)

func (i *Interpreter) LoadConsoleBuiltins() {
    vm := i.vm
    consoleObj, err := vm.New(vm.Get("Object"))
    if err != nil {
        panic(err)
    }
    vm.Set("console", consoleObj)

    consoleObj.Set("log", func(call goja.FunctionCall) goja.Value {
        li := len(call.Arguments) - 1
        for i, arg := range call.Arguments {
            if li == i {
                fmt.Printf("%s\n", arg)
            } else {
                fmt.Printf("%s ", arg)
            }
        }
        return i.vm.ToValue(nil)
    });
}
