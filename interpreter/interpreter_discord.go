package interpreter

import (
    "github.com/dop251/goja"
    "io"
)



func (i *Interpreter) LoadDiscordBuiltins() {
    vm := i.vm
    dis, err := vm.New(vm.Get("Object"))
    if err != nil {
        panic(err)
    }
    vm.Set("discord", dis)


    dis.Set("getMessage", func(call goja.FunctionCall) goja.Value {
        return i.Receive(call)
    });

    dis.Set("sendMessage", func(data string) goja.Value {
        return i.Send(data)
    });

    dis.Set("uploadFile", func(name string, contentType string, r io.Reader) goja.Value {
        return i.SendFile(name, contentType, r)
    });
    
}

func (i *Interpreter) Receive(call goja.FunctionCall) goja.Value {
    return i.vm.ToValue(i.rw.Read())
}

func (i *Interpreter) Send(data string) goja.Value {
    i.rw.WriteString(data)
    return i.vm.ToValue(nil)
}


func (i *Interpreter) SendFile(name string,contentType string, r io.Reader) goja.Value {
    i.rw.WriteFile(name, contentType, r)
    return i.vm.ToValue(nil)
}

