package main
import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"strings"
	"archive/zip"
	"os"
)



func CheckAttachments(m []*discordgo.MessageAttachment) (string, bool) {
	if len(m) == 0 {
		return "Must add at least one attachment", false
	}

	if len(m) > 1 {
		return "Must have only one attachment", false
	}
	za := m[0]
	if !strings.HasSuffix(za.URL, ".zip") {
		return "Attachment must be a zip", false
	}
	
	return "", true
}

func CheckZipEntries(archive *zip.ReadCloser) (string, bool) {
	var saveName string
	var fileIdx int
	// Get save name
	for idx, file := range archive.File {
		if !strings.HasSuffix(file.Name, ".bin") {
			saveName = file.Name
			fileIdx = idx
		}
	}
	if archive.File[(fileIdx + 1)%2].Name != saveName + ".bin" {
		return "Mismatch save files detected. Unexpected differences in file names.", false
	}

	for _, file := range archive.File {
		if strings.HasSuffix(file.Name, ".bin") {
			if file.UncompressedSize64 != 96 {
				return "Invalid .bin file size.", false
			}
		} else {
			const saveBlocks = 1 << 15
			if file.UncompressedSize64 % saveBlocks != 0 {
				return "Save image has incorrect size. Likely corrupted?", false
			}
			sizeRange := file.UncompressedSize64 >> 15
			const sizeMin = 96
			const sizeMax = 1 << 15
			if sizeRange < sizeMin {
				return "Save image is too small", false
			} else if sizeRange > sizeMax {
				return "Save image is too large", false
			}
		}
	}
	return "", true
}
func CheckZip(zp string) (string, bool) {
	archive, err := zip.OpenReader(zp)
	if err != nil {
		fmt.Println(err)
		return "There was an issue opening the file", false
	}
	defer archive.Close()
	// Count files
	if len(archive.File) > 2 {
		return "Too many files in zip. Only 2 files allowed.", false
	} else if len(archive.File) < 2 {
		return "Too little files in zip. Must have 2 files.", false
	}

	if resp, ok := CheckZipEntries(archive); !ok {
		return resp, false
	}
	return "", true
}

func DoOnline(s *discordgo.Session, m*discordgo.MessageCreate) {
	if resp, ok := CheckAttachments(m.Attachments); !ok {
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
	if resp, ok := CheckZip(zn); !ok {
		s.ChannelMessageSend(m.ChannelID, resp)
		return
	}
	s.ChannelMessageSend(m.ChannelID, "Zip passed all checks.")
}

func CommandHandler(req <-chan ClientRequest, resp chan<- string) {
	// var cs ClientState
	for {
		cr := <- req
		s := cr.Session
		m := cr.Message
		clientId := m.Author.ID
		if m.Content == "online" {
			DoOnline(s,m)
		}
		resp <- "end:" + clientId	
		break
	}	
}
