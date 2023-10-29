package main
import (
	"strings"
	"archive/zip"
)
func CheckSaveZip(archive *zip.ReadCloser) (string, bool) {
	// Count files
	fileCount := 0
	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}
		fileCount += 1
	}

	if fileCount > 2 {
		return "Too many files in zip. Only 2 files allowed.", false
	} else if fileCount < 2 {
		return "Too little files in zip. Must have 2 files.", false
	}

	if resp, ok := CheckSaveZipEntries(archive); !ok {
		return resp, false
	}
	return "", true
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
