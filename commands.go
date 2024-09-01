package main

import (
    "strings"
    "github.com/bwmarrin/discordgo"
    "github.com/Team-Alua/kat/userfs"
    "github.com/Team-Alua/kat/umountfs"
    "github.com/Team-Alua/kat/interpreter"
    "github.com/Team-Alua/kat/discord"
    "github.com/dop251/goja"
)


func InterpreterLoopMain(fn string, fs *umountfs.UmountFS, rw discord.ReadWriter) {
    for true {
        interp := interpreter.NewInterpreter(rw, fs)
        ie := interp.Run(fn)
        if gie, ok := ie.(*goja.InterruptedError); ok{

            cmd, ok := gie.Value().(string)
            if ok {
                if strings.HasPrefix(cmd, "run") {
                    fn = strings.Trim(cmd[3:], " ")
                    continue
                }

            } else {
                err := gie.Value().(error)
                rw.WriteString(err.Error())
            }
        } else if ie != nil {
            rw.WriteString(ie.Error())
        }
        break;
    }

}

func InterpreterLoop(req <-chan discord.ClientRequest, resp chan<- string, s *discordgo.Session, m *discordgo.MessageCreate) {
    rw := discord.NewReadWriter(s, req, m.ChannelID)
    fn := "default"
    mfs, err := userfs.Create(m.Author.ID)
    if err != nil {
        rw.WriteString(err.Error())
        return
    }
    InterpreterLoopMain(fn, mfs, rw)
    resp <- m.ChannelID + "_" + m.Author.ID
}


