package command

import (
	"bufio"
	"main/components/session"
	"strings"
)

// RequestUploadCommand responding for construction "REQUEST_UPLOAD" command
type RequestUploadCommand struct {
	Cmd string
}

// Row serializes command
func (c RequestUploadCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c RequestUploadCommand) Process(ctx session.Session) error {
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

	for _, upload := range split {
		if upload == "" {
			continue
		}
		cmd := CreateDownloadCommand(upload)
		_, err := ctx.GetConn().Write(cmd.Row())
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateRequestUploadCommand constructs RequestUploadCommand
func CreateRequestUploadCommand() Command {
	return &RequestUploadCommand{Cmd: "REQUEST_UPLOAD"}
}
