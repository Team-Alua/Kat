package main
import (
	"bufio"
	"encoding/json"
	"net"
	"fmt"
	"time"
)

type SaveClient struct {
	ip string
	port string
	reader *bufio.Reader
	conn net.Conn
}

func NewSaveClient(ip, port string) *SaveClient {
	return &SaveClient{ip: ip, port:port}
}

func (sc *SaveClient) Connect() (string, bool) {
	conn, err := net.DialTimeout("tcp", sc.ip + ":" + sc.port, 5 * time.Second)
	if err != nil {
		return "Failed to connect to save server.", false
	}
	sc.conn = conn
	sc.reader = bufio.NewReader(conn)
	return "", true
}

func (sc *SaveClient) Send(wb []byte) error {
	o := 0
	l := len(wb)
	for o < l {
		n, err := sc.conn.Write(wb[o:])
		if err != nil {
			return err
		}
		o += n
	} 
	return nil
}

type SaveServerResponse  struct {
	ResponseType string `json:"ResponseType"`
	Code string `json:"code"`
	Json string `json:"json"`
}

func (sc *SaveClient) Dump(saveName string, dumpFolder string) (string, bool) {
	request := fmt.Sprintf(`{"RequestType": "rtDumpSave", "sourceSaveName": "%s","targetFolder": "%s", "dumpOnly": []}` + "\n", saveName, dumpFolder)
	if err := sc.Send([]byte(request)); err != nil {
		return "Failed to send dump command", false
	}

	line, err := sc.reader.ReadSlice(byte('\n'))
	if err != nil {
		return "Failed to dump save", false
	}
	data :=	SaveServerResponse{}
	json.Unmarshal(line, &data)
	if data.ResponseType == "srInvalid" {
		return data.Code, false
	}
	return "", true
	
}

func (sc *SaveClient) Update(saveName string, updateFolder string) (string, bool) {
	request := fmt.Sprintf(`{"RequestType": "rtUpdateSave", "targetSaveName": "%s","sourceFolder": "%s", "selectOnly": []}` + "\n", saveName, updateFolder)
	if err := sc.Send([]byte(request)); err != nil {
		return "Failed to send update command", false
	}

	line, err := sc.reader.ReadSlice(byte('\n'))
	if err != nil {
		return "Failed to update save", false
	}

	data :=	SaveServerResponse{}
	json.Unmarshal(line, &data)
	if data.ResponseType == "srInvalid" {
		return data.Code, false
	}
	return "", true
}

type SaveListEntry struct {
	Kind string `json:"kind"`
	Path string `json:"path"`
	Size uint64 `json:"size"`
	Mode uint64 `json:"mode"`
	Uid  int64 `json:"uid"`
	Gid  int64 `json:"gid"`
}

func (sc *SaveClient) List(saveName string) ([]*SaveListEntry, string, bool) {
	var listEntries []*SaveListEntry
	request := fmt.Sprintf(`{"RequestType": "rtListSaveFiles", "listTargetSaveName": "%s"}` + "\n", saveName)
	if err := sc.Send([]byte(request)); err != nil {
		return listEntries, "Failed to send list command", false
	}
	line, err := sc.reader.ReadSlice(byte('\n'))
	if err != nil {
		return listEntries, "Failed to read list response", false
	}
	// Check if the response is a not okay
	data :=	SaveServerResponse{}
	json.Unmarshal(line, &data)
	if data.ResponseType == "srInvalid" {
		return listEntries, data.Code, false
	} else if data.ResponseType == "srJson" {
		json.Unmarshal([]byte(data.Json), &listEntries)
	}
	return listEntries, "", true
}

func (sc *SaveClient) Disconnect() error {
	return sc.conn.Close()
}

