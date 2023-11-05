package main
import (
	"bufio"
	"encoding/json"
	"net"
	"time"
)


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

func (sc *SaveClient) Send(req *SaveClientRequest) error {
	wb, err := json.Marshal(req)
	if err != nil {
		return err
	}
	wb = append(wb, byte('\n'))
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

func (sc *SaveClient) Dump(saveName string, dumpFolder string) (string, bool) {
	dcr := DumpClientRequest{SaveName: saveName, TargetFolder: dumpFolder, SelectOnly: make([]string, 0)}
	request := SaveClientRequest{RequestType: "rtDumpSave", Dump: &dcr}
	if err := sc.Send(&request); err != nil {
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
	ucr := UpdateClientRequest{SaveName: saveName, SourceFolder: updateFolder, SelectOnly: make([]string, 0)}
	request := SaveClientRequest{RequestType: "rtUpdateSave", Update: &ucr}
	if err := sc.Send(&request); err != nil {
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


func (sc *SaveClient) List(saveName string) ([]*SaveListEntry, string, bool) {
	var listEntries []*SaveListEntry
	lcr := ListClientRequest{SaveName: saveName}
	request := SaveClientRequest{RequestType: "rtUpdateSave", List: &lcr}
	if err := sc.Send(&request); err != nil {
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

