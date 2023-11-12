package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"github.com/Team-Alua/kat/userfs"
)

func getScript(fn string) (string, error) {
	body, err := ioutil.ReadFile(fn + ".js")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func InterpreterLoop(req <-chan ClientRequest, resp chan<- string, s *discordgo.Session, m *discordgo.MessageCreate) {
	rw := NewDiscordReadWriter(s, req, m.ChannelID)
	fn := "default"
	mfs, err := userfs.Create(m.Author.ID)
	if err != nil {
		rw.WriteString(err.Error())
		return
	}
	for true {
		code, err := getScript(fn)
		if err != nil {
			code = fmt.Sprintf(`
				send("There was an error opening %s");
			`, fn)
		}
		interp := NewInterpreter(rw, mfs)
		ie := interp.Run(fn, code)
		
		if gie, ok := ie.(*goja.InterruptedError); ok{
			cmd := gie.Value().(string)
			if strings.HasPrefix(cmd, "run") {
				fn = strings.Trim(cmd[3:], " ")
				continue
			}
		} else if ie != nil {
			rw.WriteString(ie.Error())
		}
		break;
	}
	resp <- m.ChannelID + "_" + m.Author.ID
}
