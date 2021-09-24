package command

import (
	"main/components/session"
)

// TimeCommand responding for construction "TIME" command
type TimeCommand struct {
	Cmd string
}

// Row serializes command
func (c TimeCommand) Row() []byte {
	return []byte(c.Cmd + "\n")
}

// Process writes serialized command to server
func (c TimeCommand) Process(ctx session.Session) error {
	_, err := ctx.GetConn().Write(c.Row())
	if err != nil {
		return err
	}
	return nil
}

// CreateTimeCommand constructs TimeCommand with TimeCommand.Cmd == "TIME"
func CreateTimeCommand() Command {
	return &TimeCommand{Cmd: "TIME"}
}
