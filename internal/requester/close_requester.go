package requester

import (
	"SSoC/internal/session"
)

// CloseRequester responding for construction "CLOSE" command
type CloseRequester struct {
	Cmd string
}

// Row serializes command "CLOSE"
func (c CloseRequester) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c CloseRequester) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateCloseRequester constructs CloseRequester with CloseRequester.Cmd == "CLOSE"
func CreateCloseRequester() Requester {
	return &CloseRequester{Cmd: "CLOSE"}
}
