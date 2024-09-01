package interpreter

import (
    "github.com/dop251/goja"
    "unicode/utf8"
)

func (i *Interpreter) LoadTextBuiltins() {
    vm := i.vm
    vm.Set("decodeText", func(call goja.FunctionCall) goja.Value {
        li := len(call.Arguments)
        if li < 1 {
            panic(vm.ToValue("Invalid argument count."))
        }
        bufArg := call.Argument(0).Export()
        buf, ok := bufArg.([]byte)
        if !ok {
            panic(vm.ToValue("Must be array value!"))
        }
        if !utf8.Valid(buf) {
            panic(vm.ToValue("Invalid utf8 string."))
        }
        return i.vm.ToValue(string(buf))
    });

    vm.Set("encodeText", func(call goja.FunctionCall) goja.Value {
        li := len(call.Arguments)
        if li < 1 {
            panic(vm.ToValue("Invalid argument count."))
        }
        strArg := call.Argument(0).Export()
        str, ok := strArg.(string)
        if !ok {
            panic(vm.ToValue("Must be a string value!"))
        }

        if !utf8.ValidString(str) {
            panic(vm.ToValue("Invalid utf8 string."))
        }
        buffer := vm.ToValue(vm.NewArrayBuffer([]byte(str)))
        arr, _ := vm.New(vm.Get("Uint8Array"), buffer)
        return arr;
    });
}
