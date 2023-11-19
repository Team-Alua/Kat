package discord

import (
	"github.com/bwmarrin/discordgo"
)

type ClientRequest struct {
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

func NewClientRequest(s *discordgo.Session, m *discordgo.MessageCreate) ClientRequest {
	return ClientRequest{Session: s, Message: m}
}
