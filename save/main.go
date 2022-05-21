package main

import (
	"io"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"bytes"
	"io/ioutil"
	"log"
	"strings"
	"net/http"
	"time"
	"path"
	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	httpClient = getHttpClient()
	Token string
)

func getHttpClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Print("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content != "" {
		log.Print("Message content: ", m.Content)
	}

	if len(m.Attachments) > 0 {
		processAttachments(m.Attachments)
	}
}

// https://golangcode.com/download-a-file-from-a-url/
// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// https://github.com/NyanKiyoshi/deletednt-discord
func processAttachment(attachment *discordgo.MessageAttachment) {
	// Get the attachment using a HEAD request and ensure it succeed: no error and HTTP 200
	if response, err := httpClient.Head(attachment.URL); err != nil {
		log.Fatalf("failed to get %s: %s", attachment.URL, err)
	} else if response.StatusCode != 200 {
		log.Fatalf(
			"failed to get %s, got response code: %d",
			attachment.URL, response.StatusCode)
	} else {
		log.Print("successfully pre-fetched ", attachment.URL)
		url, file := path.Split(attachment.URL)
		log.Print("base url: ", url)
		log.Print("file: ",  file)
		err := DownloadFile(file, attachment.URL)
		if err != nil {
			log.Fatalf("Error retrving file to local disk! %s", err)
		}
			log.Print("Downloaded: " + attachment.URL)
	}
}

func processAttachments(attachments []*discordgo.MessageAttachment) {
	for _, attachment := range attachments {
		processAttachment(attachment)
	}
}

func retrieveDeletedAttachment(attachment *discordgo.MessageAttachment) *discordgo.File {
	response, err := httpClient.Get(attachment.URL)
	if err == nil {
		defer response.Body.Close()

		if response.StatusCode == 200 {
			var buffer []byte
			buffer, err := ioutil.ReadAll(response.Body)

			if err == nil {
				reader := bytes.NewReader(buffer)
				contentType := strings.Split(response.Header.Get("Content-Type"), ";")[0]

				return &discordgo.File{
					Name:        attachment.Filename,
					ContentType: contentType,
					Reader:      reader,
				}
			}
		}
	}

	if err == nil {
		err = fmt.Errorf("received unexpected status %d", response.StatusCode)
	}
	log.Printf("failed to get %s: %s", attachment.URL, err)
	return nil
}
