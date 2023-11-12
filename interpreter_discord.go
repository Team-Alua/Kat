package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dop251/goja"
	"io"
)

type DiscordReadWriter struct {
	session *discordgo.Session
	channelID string
	in <-chan ClientRequest
}

func NewDiscordReadWriter(s *discordgo.Session, 
						  in <-chan ClientRequest, 
						  channelID string) *DiscordReadWriter {

	rw := &DiscordReadWriter{}
	rw.session = s
	rw.channelID = channelID
	rw.in = in
	return rw
}

func (rw *DiscordReadWriter) WriteString(msg string) {
	s := rw.session
	chanID := rw.channelID
	s.ChannelMessageSend(chanID, msg)
}

func (rw *DiscordReadWriter) WriteFile(name string, contentType string, r io.Reader) {
	s := rw.session
	chanID := rw.channelID
	data := &discordgo.MessageSend{Files: []*discordgo.File{{Name: name,ContentType: contentType,Reader: r}}}
	s.ChannelMessageSendComplex(chanID, data)
}

func (rw *DiscordReadWriter) Read() *discordgo.MessageCreate {
	cr := <- rw.in
	// Update session just in case
	rw.session = cr.Session
	return cr.Message
}


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
