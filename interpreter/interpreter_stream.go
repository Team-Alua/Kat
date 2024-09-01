package interpreter

import (
    "github.com/dop251/goja"
    "io"
    "bufio"
)


func (i *Interpreter) LoadStreamBuiltins() {
    vm := i.vm
    vm.Set("StreamWriter", func(call goja.ConstructorCall) *goja.Object {
        fp, ok := call.Argument(0).Export().(io.Writer)
        if !ok {
            panic("First argument must be Readable.")
        }
        writer := bufio.NewWriter(fp)

        call.This.Set("writeString", func(icall goja.FunctionCall) goja.Value {
            // flush
            str := call.Argument(0).Export().(string)
            if !ok {
                panic("First argument must be a string.")
            }
            n, err := writer.WriteString(str)
            if err != nil {
                panic(err)
            }
            return vm.ToValue(n)
        })

        call.This.Set("writeLine", func(icall goja.FunctionCall) goja.Value {
            // write string and add newline
            // flush
            str := call.Argument(0).Export().(string)
            if !ok {
                panic("First argument must be a string.")
            }
            n, err := writer.WriteString(str + "\n")
            if err != nil {
                panic(err)
            }
            return vm.ToValue(n)
        })

        call.This.Set("write", func(icall goja.FunctionCall) goja.Value {
            // TypedArray => ArrayBuffer
            // flush
            buf, ok := call.Argument(0).Export().(goja.ArrayBuffer)
            if !ok {
                panic("Not an arraybuffer")
            }
            n, err := writer.Write(buf.Bytes())
            if err != nil {
                panic(err)
            }
            return vm.ToValue(n)
        })

        call.This.Set("close", func(icall goja.FunctionCall) goja.Value {
            if err := writer.Flush(); err != nil {
                panic(err)
            }

            return vm.ToValue(nil)
        })

        return nil
    })

    vm.Set("StreamReader", func(call goja.ConstructorCall) *goja.Object {
        fp, ok := call.Argument(0).Export().(io.Reader)
        if !ok {
            panic("First argument must be Readable.")
        }

        reader := bufio.NewReader(fp)

        call.This.Set("readUntil", func(icall goja.FunctionCall) goja.Value {
            d := icall.Argument(0).String();

            if len(d) != 1 {
                panic("First argument must a single character.")
            }
            b := []byte(d)
            buf, err := reader.ReadBytes(b[0])
            if err != nil {
                panic(err)
            }
            return vm.ToValue(buf)
        })

        call.This.Set("read", func(icall goja.FunctionCall) goja.Value {
            // int32
            c := icall.Argument(0).ToInteger()

            buf := make([]byte, c)
            n, err := reader.Read(buf)
            if err != nil {
                panic(err)
            }
            buf = buf[0:n]
            return vm.ToValue(buf)
        })

        call.This.Set("readLine", func(icall goja.FunctionCall) goja.Value {
            buf, err := reader.ReadBytes('\n')
            if err != nil {
                panic(err)
            }
            return vm.ToValue(string(buf))
        })

        return nil
    })
}

