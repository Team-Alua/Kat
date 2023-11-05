package main

import (
	"github.com/google/uuid"
	"github.com/bwmarrin/discordgo"
	"archive/zip"
	"os"
	"fmt"
	"strings"
	"strconv"
)

func DoResign(s *discordgo.Session, m *discordgo.MessageCreate, pzn string) {
	// Must have at least one upload
	if (pzn == "") {
		s.ChannelMessageSend(m.ChannelID, "Must have at least one upload.")
		return
	}


	var accountId uint64
	// Check message for accountId

	// Split command
	cmdPieces := strings.Fields(m.Content)

	if len(cmdPieces) != 2 {
		s.ChannelMessageSend(m.ChannelID, "Expected something like `resign 0x1234567812345678`.")
		return 

	}
	accountId, err := strconv.ParseUint(cmdPieces[1], 0, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "The supplied PSN ID was invalid.")
		return 
	}

	psarc, err := zip.OpenReader(pzn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue opening up the save zip file.")
		return 
	}
	defer psarc.Close()


	// Generate ID 
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	s.ChannelMessageSend(m.ChannelID, "Save ID:" + id)

	fc := NewFtpClient("10.0.0.5", "2121")
	if resp, ok := fc.Login(); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return 
	}
	defer fc.Kill()

	// Now upload the files to the PS4
	if resp, ok := fc.UploadSave(psarc, id); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	sc := NewSaveClient("10.0.0.5", "1234")
	if resp, ok := sc.Connect(); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return 
	}
	defer sc.Disconnect()

	s.ChannelMessageSend(m.ChannelID, "Connected to PS4 save server.")

	if resp, ok := sc.Resign(id, accountId); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}


	// Look for the shortest save name
	saveName := ""
	for _, f := range psarc.File {
		// Ignore folders
		if f.FileInfo().IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name, ".bin") {
			saveName = f.Name
		}
	}

	outzn := m.Author.ID + "_out.zip"
	outarc, err := os.Create(outzn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue creating new.zip.")
		return
	}
	defer os.Remove(outzn)

	defer outarc.Close()

	outw := zip.NewWriter(outarc)
	err = fc.ZipSave(outw, saveName, id)

	outw.Close()

	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "There was an issue updating save zip.")
		return
	}
	outarc.Seek(0, 0)
	data := &discordgo.MessageSend{Files: []*discordgo.File{{Name: "new.zip",ContentType: "application/zip",Reader: outarc}}}
	s.ChannelMessageSendComplex(m.ChannelID, data)
}
