package discord

import (
    "io"
    "github.com/bwmarrin/discordgo"
)

type ReadWriter interface {
    WriteString(msg string)
    WriteFile(name string, contentType string, r io.Reader)
    Read() *discordgo.MessageCreate
}

