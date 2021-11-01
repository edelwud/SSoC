package command

import (
	"bufio"
	"main/components/session"
	"strings"
)

// RequestDownloadCommand responding for construction "REQUEST_DOWNLOAD" command
type RequestDownloadCommand struct {
	Cmd string
}

// Row serializes command
func (c RequestDownloadCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c RequestDownloadCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}

	buf, _, err := bufio.NewReader(ctx.GetConn()).ReadLine()
	if err != nil {
		return err
	}

	result := strings.Trim(string(buf), "\n")
	split := strings.Split(result, ",")

	for _, download := range split {
		if download == "" {
			continue
		}
		cmd := CreateDownloadCommand(download)
		_, err := ctx.GetConn().Write(cmd.Row())
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateRequestDownloadCommand constructs RequestDownloadCommand
func CreateRequestDownloadCommand() Command {
	return &RequestDownloadCommand{Cmd: "REQUEST_DOWNLOAD"}
}
