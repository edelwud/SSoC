package requester

import (
	"SSoC/internal/session"
	"bufio"
	"strings"
)

// RequestUploadRequester responding for construction "REQUEST_UPLOAD" command
type RequestUploadRequester struct {
	Cmd string
}

// Row serializes command
func (c RequestUploadRequester) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c RequestUploadRequester) Process(ctx session.Session) error {
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
		cmd := CreateUploadRequester(upload, "")
		_, err := ctx.GetConn().Write(cmd.Row())
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateRequestUploadRequester constructs RequestUploadRequester
func CreateRequestUploadRequester() Requester {
	return &RequestUploadRequester{Cmd: "REQUEST_UPLOAD"}
}
