package requester

import (
	"SSoC/internal/session"
)

// EchoRequester responding for construction "ECHO <client text>" command
type EchoRequester struct {
	Cmd  string
	Text string
}

// Row serializes command
func (c EchoRequester) Row() []byte {
	return []byte(c.Cmd + " " + c.Text + "\n")
}

// Process writes serialized command to server
func (c EchoRequester) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateEchoRequester constructs EchoRequester with EchoRequester.Cmd == "ECHO"
func CreateEchoRequester(text string) Requester {
	return &EchoRequester{Cmd: "ECHO", Text: text}
}
