package main

import (
	"github.com/google/uuid"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"strings"
	"fmt"
	"os"
	"archive/zip"
)

func CheckUpdateAttachments(ma []*discordgo.MessageAttachment) (string, bool) {
	if len(ma) != 2 {
		return "Must have exactly two attachment", false
	}

	for _, m := range ma {
		u, _ := url.Parse(m.URL)
		if !strings.HasSuffix(u.Path, ".zip") {
			return "All attachments must be a zip", false
		}
	}

	return "", true
}

func CheckUpdateAttachmentNames(ma []*discordgo.MessageAttachment) (string, bool) {
	uc := 0
	for _, m := range ma {
		u, _ := url.Parse(m.URL)
		if strings.HasSuffix(u.Path, "/update.zip") {
			uc += 1
		}
	}
	if uc != 1 {
		return "There must be exactly one attachment named update.zip", false
	}

	return "", true
}

func DownloadUpdateAttachments(ma []*discordgo.MessageAttachment, prefix string) (string, bool) {
	for _, m := range ma {
		
		u, _ := url.Parse(m.URL)
		zn := "update.zip"
		if !strings.HasSuffix(u.Path, "/update.zip") {
			zn = "PS4.zip"
		}
		zn = prefix + "_" + zn

		if err := DownloadFile(zn, m.URL); err != nil {
			fmt.Println(err)
			return fmt.Sprintf("Failed to download %s.", m.URL), false
		}
	}
	return "", true
}

func DoUpdate(s *discordgo.Session, m*discordgo.MessageCreate) {
	if resp, ok := CheckUpdateAttachments(m.Attachments); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	if resp, ok := CheckUpdateAttachmentNames(m.Attachments); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	
	if resp, ok := DownloadUpdateAttachments(m.Attachments, m.Author.ID); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	pzn := m.Author.ID + "_PS4.zip"
	uzn := m.Author.ID + "_update.zip" 
	// Clean up zip files
	defer os.Remove(pzn)
	defer os.Remove(uzn)

	psarc, err := zip.OpenReader(pzn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue opening up the save zip file.")
		return 
	}
	defer psarc.Close()

	if resp, ok := CheckSaveZip(psarc); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	uparc, err := zip.OpenReader(uzn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue opening up the update zip file.")
		return 
	}
	defer uparc.Close()

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

	// Create the update folder as well
	resp, ok := fc.CreateTempFolder(id)
	if !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	tmpFolder := resp

	defer fc.DeleteStage(id)

	if resp, ok := fc.UploadDump(uparc, tmpFolder); !ok {
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

	if resp, ok := sc.Update(id, tmpFolder); !ok {
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
