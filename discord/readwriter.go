package discord

import (
    "github.com/bwmarrin/discordgo"
    "io"
)

type DiscordReadWriter struct {
    session *discordgo.Session
    channelID string
    in <-chan ClientRequest
}

func NewReadWriter(s *discordgo.Session, 
                          in <-chan ClientRequest,
                          channelID string) ReadWriter {

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

