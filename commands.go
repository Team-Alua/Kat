package main
import (
	"strings"
	"archive/zip"
)
func CheckSaveZip(archive *zip.ReadCloser) (string, bool) {
	// Count files
	if len(archive.File) > 2 {
		return "Too many files in zip. Only 2 files allowed.", false
	} else if len(archive.File) < 2 {
		return "Too little files in zip. Must have 2 files.", false
	}

	if resp, ok := CheckSaveZipEntries(archive); !ok {
		return resp, false
	}
	return "", true
}


func CheckSaveZipEntries(archive *zip.ReadCloser) (string, bool) {
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

func CommandHandler(req <-chan ClientRequest, resp chan<- string) {
	// var cs ClientState
	for {
		cr := <- req
		s := cr.Session
		m := cr.Message
		clientId := m.Author.ID
		if m.Content == "dump" {
			DoDump(s,m)
		} else if m.Content == "update" {
			DoUpdate(s,m)
		}

		resp <- "end:" + clientId	
		break
	}	
}
