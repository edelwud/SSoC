package requester

import (
	"SSoC/internal/session"
	"bufio"
	"strings"
)

// RequestDownloadRequester responding for construction "REQUEST_DOWNLOAD" command
type RequestDownloadRequester struct {
	Cmd string
}

// Row serializes command
func (c RequestDownloadRequester) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c RequestDownloadRequester) Process(ctx session.Session) error {
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
		cmd := CreateDownloadRequester(download)
		_, err := ctx.GetConn().Write(cmd.Row())
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateRequestDownloadRequester constructs RequestDownloadRequester
func CreateRequestDownloadRequester() Requester {
	return &RequestDownloadRequester{Cmd: "REQUEST_DOWNLOAD"}
}
