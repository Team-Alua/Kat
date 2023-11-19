package interpreter

import (
	"github.com/dop251/goja"
	"io"
	"bufio"
)


func (i *Interpreter) LoadStreamBuiltins() {
	vm := i.vm

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
		
		call.This.Set("readline", func(icall goja.FunctionCall) goja.Value {
			buf, err := reader.ReadBytes('\n')
			if err != nil {
				panic(err)
			}
			return vm.ToValue(string(buf))
		})

		return nil
	})
}

