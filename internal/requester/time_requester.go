package requester

import (
	"SSoC/internal/session"
)

// TimeRequester responding for construction "TIME" command
type TimeRequester struct {
	Cmd string
}

// Row serializes command
func (c TimeRequester) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c TimeRequester) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateTimeRequester constructs TimeRequester with TimeRequester.Cmd == "TIME"
func CreateTimeRequester() Requester {
	return &TimeRequester{Cmd: "TIME"}
}
