package main

import "bufio"
import "net"

type SaveClient struct {
	ip string
	port string
	reader *bufio.Reader
	conn net.Conn
}

type DumpClientRequest struct {
	SaveName string `json:"saveName"`
	TargetFolder string `json:"targetFolder"`
	SelectOnly []string `json:"selectOnly"`
}

type UpdateClientRequest struct {
	SaveName string `json:"saveName"`
	SourceFolder string `json:"sourceFolder"`
	SelectOnly []string `json:"selectOnly"`
}

type ResignClientRequest struct {
	SaveName string `json:"saveName"`
	AccountId uint64 `json:"accountId"`
}

type ListClientRequest struct {
	SaveName string `json:"saveName"`
}

type SaveClientRequest struct {
	RequestType string `json:"RequestType"`
	Dump *DumpClientRequest `json:"dump,omitempty"`
	Update *UpdateClientRequest `json:"update,omitempty"`
	Resign *ResignClientRequest `json:"resign,omitempty"`
	List *ListClientRequest `json:"list,omitempty"`
}

type SaveServerResponse  struct {
	ResponseType string `json:"ResponseType"`
	Code string `json:"code"`
	Json string `json:"json"`
}

type SaveListEntry struct {
	Kind string `json:"kind"`
	Path string `json:"path"`
	Size uint64 `json:"size"`
	Mode uint64 `json:"mode"`
	Uid  int64 `json:"uid"`
	Gid  int64 `json:"gid"`
}

