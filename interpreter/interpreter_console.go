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
            obj := arg.ToObject(vm)
            str, _ := obj.MarshalJSON()
            if li == i {
                fmt.Printf("%s\n", str)
            } else {
                fmt.Printf("%s ", str)
            }
        }
        return i.vm.ToValue(nil)
    });
}
