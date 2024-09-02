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

    vm.Set("fnv1a32", func(call goja.FunctionCall) goja.Value {
        strArg := call.Argument(0).Export()
        str, _ := strArg.(string)
        OFFSET_BASIS := uint32(0x811c9dc5)
        hash := uint32(OFFSET_BASIS)
        for i := 0; i < len(str); i++ {
            hash = hash ^ uint32(str[i])
            hash *= uint32(0x01000193)
        }
        return i.vm.ToValue(int64(hash))
    })

    // (target, offset, source)
    vm.Set("copyBuffer", func(call goja.FunctionCall) goja.Value {
        buf1 := call.Argument(0).Export().([]byte)
        offset := call.Argument(1).ToInteger()
        buf2 := call.Argument(2).Export().([]byte)
        buf2len := int64(len(buf2))
        for idx := int64(0); idx < buf2len; idx++ {
            buf1[offset + idx] = buf2[idx]
        }
        return i.vm.ToValue(nil)
    })

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
