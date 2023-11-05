package main


import (
	"github.com/bwmarrin/discordgo"
	"archive/zip"
//	"os"
	"fmt"
	"errors"
	"strings"
	"path"
	"net/url"
)


func CheckSaveZip(zippath string) error {
	archive, err := zip.OpenReader(zippath)
	if err != nil {
		fmt.Println(err)
		return errors.New("There was an issue opening up the uploaded zip file.")
	}
	defer archive.Close()

	// Count files
	fileCount := 0
	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}
		fileCount += 1
	}

	if fileCount > 2 {
		return errors.New("Too many files in zip. Only 2 files allowed.")
	} else if fileCount < 2 {
		return errors.New("Too little files in zip. Must have 2 files.")
	}

	if resp, ok := CheckSaveZipEntries(archive); !ok {
		return errors.New(resp)
	}
	return nil
}


func CheckSaveZipEntries(archive *zip.ReadCloser) (string, bool) {
	var saveName string
	var saveIdx int
	var binName string
	var binIdx int
	// Get save name
	for idx, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name, ".bin") {
			binName = file.Name
			binIdx = idx
		} else {
			saveName = file.Name
			saveIdx = idx
		} 
	}
	if binName != saveName + ".bin" {
		return "Mismatch save files detected. Unexpected differences in file names.", false
	}
	binFile := archive.File[binIdx]
	if binFile.UncompressedSize64 != 96 {
		return "Invalid .bin file size.", false
	}
	saveFile := archive.File[saveIdx]
	const saveBlocks = 1 << 15
	if saveFile.UncompressedSize64 % saveBlocks != 0 {
		return "Save image has incorrect size. Likely corrupted?", false
	}
	sizeRange := saveFile.UncompressedSize64 >> 15
	const sizeMin = 96
	const sizeMax = 1 << 15
	if sizeRange < sizeMin {
		return "Save image is too small", false
	} else if sizeRange > sizeMax {
		return "Save image is too large", false
	}

	return "", true
}

func GetUploadZipName(authorId string, attachmentId string) string {
	return authorId + "_" + attachmentId + ".zip"
}

func CheckAndDownloadAttachments(ma []*discordgo.MessageAttachment, authorId string) ([]string, bool) {
	errs := make([]string, len(ma))
	success := true
	downloaded := 0
	for _, m := range ma {
		var err error
		zn := GetUploadZipName(authorId, m.ID)
		if !strings.HasSuffix(m.Filename, ".zip") {
			errMsg := fmt.Sprintf("%s is not a zip.", m.Filename)
			errs = append(errs, errMsg)
		} else if err = DownloadFile(zn, m.URL); err != nil {
			errMsg := fmt.Sprintf("Failed to download %s.", m.Filename)
			errs = append(errs, errMsg)
		} else if err = CheckSaveZip(zn); err != nil {
			fmt.Println(err)
			errMsg := fmt.Sprintf("%s : %s.", m.Filename, err.Error())
			errs = append(errs, errMsg)
		}
		if err != nil {
			fmt.Println(m.URL, zn, err)
			os.Remove(zn)
		} else {
			downloaded += 1
		}
	}

	if (len(errs) > 0) {
		success = false
		errs = append(errs, fmt.Sprintf("Downloaded %d out of %d zips.", downloaded, len(ma)))
	}

	return errs, success
}

func DoUpload(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	// Parse message for links
	// Should only have links otherwise ignore
	links := strings.Fields(m.Content)
	for _ , link := range links {
		u, err := url.ParseRequestURI(link)
		// Only want valid links
		if err != nil {
			return false
		}

		am := &discordgo.MessageAttachment{}
		am.ID = path.Base(path.Dir(u.Path))
		am.Filename = path.Base(u.Path)
		am.URL = link
		m.Attachments = append(m.Attachments, am)
	}
	// Ignore
	if len(m.Attachments) == 0 {
		return false
	}

	if resp, ok := CheckAndDownloadAttachments(m.Attachments, m.Author.ID); !ok {
		s.ChannelMessageSend(m.ChannelID, strings.Join(resp, "\n"))
		return true
	}
	s.ChannelMessageSend(m.ChannelID, "Downloaded all attachments.")
	return true
}

