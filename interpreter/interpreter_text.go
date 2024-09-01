package interpreter

import (
    "github.com/dop251/goja"
)

func (i *Interpreter) LoadTextBuiltins() {
    vm := i.vm
    vm.Set("decodeText", func(call goja.FunctionCall) goja.Value {
        li := len(call.Arguments)
        if li < 1 {
            panic("Invalid argument count.")
        }
        bufArg := call.Argument(0).Export()
        buf, ok := bufArg.([]byte)
        if !ok {
            panic("Must be array value!")
        }

        return i.vm.ToValue(string(buf))
    });
}
