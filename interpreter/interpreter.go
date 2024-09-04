package interpreter

import (
    "fmt"
    "io/ioutil"
    "github.com/dop251/goja"
    "github.com/Team-Alua/kat/umountfs"
    "github.com/Team-Alua/kat/discord"
)

func GetScript(fn string) (string, error) {
    body, err := ioutil.ReadFile("scripts/" + fn + ".js")
    if err != nil {
        return "", err
    }
    return string(body), nil
}


type Interpreter struct {
    vm *goja.Runtime
    rw discord.ReadWriter
    fs *umountfs.UmountFS
}

func NewInterpreter(rw discord.ReadWriter, fs *umountfs.UmountFS) *Interpreter {
    i := &Interpreter{} 
    i.vm = goja.New()
    i.rw = rw
    i.fs = fs
    return i
}

func (i *Interpreter) Run(name string) error {
    vm := i.vm
    i.LoadBuiltins()
    defer func() {
        i.fs.UnmountAll()
    }()
    code, err := GetScript(name)
    if err != nil {
        code = fmt.Sprintf(`
            send("There was an error opening %s");
        `, name)
    }
    _, err = vm.RunScript(name, code)
    return err
}

func (i *Interpreter) LoadBuiltins() {
    i.LoadFsBuiltins()
    i.LoadDiscordBuiltins()
    i.LoadHttpBuiltins()
    i.LoadStreamBuiltins()
    i.LoadConsoleBuiltins()
    i.LoadTextBuiltins()
    vm := i.vm



    vm.Set("run", func(name string) goja.Value {
        code, err := GetScript(name)
        if err != nil {
            vm.Interrupt(err)
            return vm.ToValue(nil)
        }
        value, err := i.vm.RunScript(name, code)
        if err != nil {
            vm.Interrupt(err)
            return vm.ToValue(nil)
        }
        return value
    });

    vm.Set("exit", func(call goja.FunctionCall) goja.Value {
        vm.Interrupt("exit")
        return vm.ToValue(nil)
    });
}

