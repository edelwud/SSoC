package command

type EchoCommand struct {
	Cmd  string
	Text string
}

func (c EchoCommand) Row() []byte {
	return []byte(c.Cmd + " " + c.Text + "\n")
}

func CreateEchoCommand(text string) Command {
	return &EchoCommand{Cmd: "ECHO", Text: text}
}
