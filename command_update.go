package main

import (
	"github.com/google/uuid"
	"github.com/bwmarrin/discordgo"
	"strings"
	"fmt"
	"os"
	"archive/zip"
)


func GetUpdateAttachmentName(ma *discordgo.MessageAttachment, authorId string) string {
	return authorId + "_" + ma.ID + "_update.zip"
}

func CheckUpdateAttachments(ma []*discordgo.MessageAttachment) (string, bool) {
	if len(ma) != 1 {
		return "Please upload a single zip.", false
	}

	if !strings.HasSuffix(ma[0].Filename, ".zip") {
		return "Upload must be a zip file", false
	}

	return "", true
}

func DownloadUpdateAttachments(m *discordgo.MessageCreate) (string, bool) {
	ma := m.Attachments[0]
	zn := GetUpdateAttachmentName(ma, m.Author.ID)
	if err := DownloadFile(zn, ma.URL); err != nil {
		fmt.Println(err)
		return fmt.Sprintf("Failed to download %s.", ma.Filename), false
	}

	return "", true
}

func DoUpdate(s *discordgo.Session, m *discordgo.MessageCreate, pzn string) {
	// Must have at least one upload
	if (pzn == "") {
		s.ChannelMessageSend(m.ChannelID, "Must have at least one upload.")
		return
	}

	if resp, ok := CheckUpdateAttachments(m.Attachments); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}

	if resp, ok := DownloadUpdateAttachments(m); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	uzn := GetUpdateAttachmentName(m.Attachments[0], m.Author.ID)

	// Clean up zip files
	defer os.Remove(uzn)

	psarc, err := zip.OpenReader(pzn)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "There was an issue opening up the save zip file.")
		return 
	}
	defer psarc.Close()

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
