package main
import (
	"io"
	"os"
	"net/http"
	"errors"
)
// https://golangcode.com/download-a-file-from-a-url/
// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadToWriter(w io.Writer, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("There was an issue with the provided url.")
	}

	// Write the body to file
	_, err = io.Copy(w, resp.Body)
	return err
}

func DownloadFile(path string, url string) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return DownloadToWriter(out, url)
}
