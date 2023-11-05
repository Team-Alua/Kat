package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"
)
func RebuildUploadList(authorId string, uploads *[]string) {
	files, err := ioutil.ReadDir(".")
    if err != nil {
        fmt.Println(err)
		return
    }

	// Clear slice
	*uploads = nil
	*uploads = make([]string, 0)
    for _, file := range files {
		if file.IsDir() {
			continue
		}
		fn := file.Name()
		if strings.HasPrefix(fn, authorId) && strings.HasSuffix(fn, "zip") {
			*uploads = append(*uploads, fn)
		} 
    }
}


func CommandHandler(req <-chan ClientRequest, resp chan<- string) {
	var uploads []string = make([]string, 0)
	var clientId string
	firstTime := true
	for {
		cr := <- req
		s := cr.Session
		m := cr.Message
		clientId = m.Author.ID
		if firstTime {
			firstTime = false
			RebuildUploadList(m.Author.ID, &uploads)
		}
		fmt.Println("Received request", m.Content)
		pszip := ""
		if len(uploads) > 0 {
			pszip = uploads[len(uploads) - 1]
		}
		if m.Content == "dump" {
			DoDump(s,m,pszip)
		} else if m.Content == "update" {
			DoUpdate(s,m,pszip)
		} else if m.Content == "end" {
			s.ChannelMessageSend(m.ChannelID, "Enjoy.")
			break
		} else if (strings.HasPrefix(m.Content, "resign")) {
			DoResign(s, m, pszip)
		} else if DoUpload(s,m) {
			RebuildUploadList(m.Author.ID, &uploads)
		}
		resp <- clientId
	}
	for _, upload := range uploads {
		os.Remove(upload)
	}
	resp <- "end:" + clientId	
}
