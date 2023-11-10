package main

import (
	"archive/zip"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
	"io"
	"path/filepath"
	"fmt"
)

type FtpClient struct {
	ip string
	port string
	conn *ftp.ServerConn
}


func NewFtpClient(ip string, port string) *FtpClient {
	return &FtpClient{ip: ip, port: port}
}

func (c *FtpClient) Login() (string, bool) {
	conn, err := ftp.Dial(c.ip + ":" + c.port, ftp.DialWithTimeout(1*time.Second))
	if err != nil {
		return "Failed to connect to ftp server.", false
	}

	err = conn.Login("anonymous", "anonymous")
	if err != nil {
		return "Failed to login to PS4.", false
	}
	c.conn = conn
	return "", true
}

func (c *FtpClient) CreateTempFolder(id string) (string, bool) {
	tmp := "/data/" + id + "_stage"
	if err := c.conn.MakeDir(tmp); err != nil {
		return "Failed to create temp folder.", false
	}
	return tmp, true
}


func (c *FtpClient) UploadSave(archive *zip.ReadCloser, id string) (string, bool) {
	errMsg := ""
	for _, f := range archive.File {
		// Do not upload folders to the PS4
		if f.FileInfo().IsDir() {
			continue
		}

		tp := "/data/" + id	
		if strings.HasSuffix(f.Name, ".bin") {
			tp += ".bin"
		}
		rc, err := f.Open()
		
		if err != nil {
			errMsg = "Failed to retrieve save file from zip."
			rc.Close()
			break
		}

		err = c.conn.Stor(tp, rc)
		if err != nil {
			errMsg = "Failed to upload save file to PS4."
			break
		}
	}
	return errMsg, errMsg == ""
}


func (c *FtpClient) ZipSave(archive *zip.Writer, saveName string, remoteSaveName string) error {
	suffix := ""
	var fError error
	for i := 1; i <= 2; i++ {
		if i == 2 {
			suffix = ".bin"
		}
		w, err := archive.Create(saveName + suffix)
		if err != nil {
			fError = err
			break
		}

		tp := "/data/" + remoteSaveName + suffix
		r, err := c.conn.Retr(tp)
		if err != nil {
			fError = err
			break
		}
		_, err = io.Copy(w, r)
		r.Close()
		if err != nil {
			fError = err
			break
		}
	}
	return fError
}

func (c *FtpClient) DeleteSave(id string) (string, bool) {
	if err := c.conn.Delete("/data/" + id); err != nil {
		return "Failed to delete save on PS4.", false
	}
	if err := c.conn.Delete("/data/" + id + ".bin"); err != nil {
		return "Failed to delete save on PS4.", false
	}
	return "", true
}


func (c *FtpClient) UploadDump(archive *zip.ReadCloser, updateFolder string) (string, bool) {
	createdDirs := make(map[string]bool)
	errMsg := ""

	for _, f := range archive.File {
		if f.FileInfo().IsDir() {
			continue
		}
		root := ""
		roots := strings.Split(f.Name, "/")
		for _, d := range roots[0:len(roots)-1] {
			if d == "" || d == ".." || d == "." {
				continue
			}
			root += d + "/"
			if !createdDirs[root] {
				createdDirs[root] = true
				c.conn.MakeDir(updateFolder + "/" + root)
			}
		}

		rc, err := f.Open()

		if err != nil {
			errMsg = "Failed to retrieve dump file from zip."
			rc.Close()
			break
		}
		err = c.conn.Stor(updateFolder + "/" + f.Name, rc)
		if err != nil {
			errMsg = "Failed to upload dump file to PS4."
			break
		}
	}

	return errMsg, errMsg == ""
}

func (c *FtpClient) ZipDump(archive *zip.Writer, dumpFolder string) error {
	walker := c.conn.Walk(dumpFolder)
	var fError error
	for walker.Next() {
		rp := walker.Path()
		if walker.Stat().Type == ftp.EntryTypeFolder {
			continue
		}
		zp, _ := filepath.Rel(dumpFolder, walker.Path())
		w, err := archive.Create(zp)
		if err != nil {
			fError = err
			break
		}
		r, err := c.conn.Retr(rp)
		if err != nil {
			fError = err
			break
		}
		_, err = io.Copy(w, r)
		r.Close()
		if err != nil {
			fError = err
			break
		}
	}
	return fError
}

func (c *FtpClient) DeleteFolder(path string) {
	c.conn.ChangeDir(filepath.Dir(path))
	c.conn.RemoveDir(filepath.Base(path))
}

func (c *FtpClient) DeleteStage(saveName string) {
	c.conn.Delete("/data/" + saveName)
	c.conn.Delete("/data/" + saveName + ".bin")
	stageFolder := "/data/" + saveName + "_stage"
	walker := c.conn.Walk(stageFolder)
	defer c.DeleteFolder(stageFolder)
	for walker.Next() {
		if walker.Err() != nil {
			fmt.Println("Failed to walk", walker.Path())
			break
		}

		if walker.Stat().Type == ftp.EntryTypeFile {
			c.conn.Delete(walker.Path())
		} else if walker.Stat().Type == ftp.EntryTypeFolder {
			defer c.DeleteFolder(walker.Path())
		} 
	}
}

func (c *FtpClient) Kill() bool {
	return c.conn.Quit() == nil
}

