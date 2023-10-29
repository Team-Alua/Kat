package main

import (
	"github.com/google/uuid"
	"github.com/bwmarrin/discordgo"
	"strings"
	"archive/zip"
	"os"
	"net/url"
	"bytes"
	"fmt"
)

func CheckDumpAttachments(ma []*discordgo.MessageAttachment) (string, bool) {
	if len(ma) == 0 {
		return "Must add at least one attachment", false
	}

	if len(ma) > 1 {
		return "Must have only one attachment", false
	}
	u, _ := url.Parse(ma[0].URL)
	if !strings.HasSuffix(u.Path, ".zip") {
		return "Attachment must be a zip", false
	}
	
	return "", true
}


func DoDump(s *discordgo.Session, m*discordgo.MessageCreate) {
	if resp, ok := CheckDumpAttachments(m.Attachments); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	zn := m.Author.ID + "_PS4.zip"
	if err := DownloadFile(zn, m.Attachments[0].URL); err != nil {
		fmt.Println(err)
		s.ChannelMessageSend(m.ChannelID, "Failed to download attachment.")
		return
	}
	defer os.Remove(zn)
	s.ChannelMessageSend(m.ChannelID, "Downloaded attachment.")

	
	archive, err := zip.OpenReader(zn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue opening the file")
		return 
	}
	defer archive.Close()
	if resp, ok := CheckSaveZip(archive); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Zip passed all checks.")
	
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
	if resp, ok := fc.UploadSave(archive, id); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	// Create the folder as well
	resp, ok := fc.CreateTempFolder(id)
	if !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	tmpFolder := resp

	sc := NewSaveClient("10.0.0.5", "1234")
	if resp, ok := sc.Connect(); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return 
	}
	defer sc.Disconnect()

	s.ChannelMessageSend(m.ChannelID, "Connected to PS4 save server.")

	if resp, ok := sc.Dump(id, tmpFolder); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	
	defer fc.DeleteStage(id)

	s.ChannelMessageSend(m.ChannelID, "Dumped save.")
	
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	
	err = fc.ZipDump(w, tmpFolder)
	w.Close()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Failed zip up dump.")
		return 
	}
	data := &discordgo.MessageSend{Files: []*discordgo.File{{Name: "dump.zip",ContentType: "application/zip",Reader: buf}}}
	s.ChannelMessageSendComplex(m.ChannelID, data)

}

