package cmdline

import (
    "bufio"
    "os"
    "io"
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/Team-Alua/kat/discord"
    "encoding/json"
)

type CmdReadWriter struct {
}

func NewReadWriter() discord.ReadWriter {
    rw := &CmdReadWriter{}
    return rw
}

func (rw *CmdReadWriter) WriteString(msg string) {
    fmt.Println(msg)
}

func (rw *CmdReadWriter) WriteFile(name string, contentType string, r io.Reader) {
    fmt.Println("Writing file ", name, " contentType: ", contentType)
}

func (rw *CmdReadWriter) Read() *discordgo.MessageCreate {
    reader := bufio.NewReader(os.Stdin)
    text, _ := reader.ReadString('\n')
    var msg discordgo.Message
    if err := json.Unmarshal([]byte(text), &msg); err != nil {
        panic(err)
    }
    msgCreate := &discordgo.MessageCreate{Message: &msg}
    return msgCreate
}

