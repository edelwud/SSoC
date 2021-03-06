package command

import (
	"errors"
	"regexp"
	"strings"
)

// Command contains necessary for command execution parameters
type Command struct {
	Cmd        string
	Parameters []string
}

// MaxParametersCount maximum number of Command.Parameters
const MaxParametersCount = 100

// ParseCommand parses client command via regex, first parameter is Command.Cmd, next are Command.Parameters
func ParseCommand(command string) (Command, error) {
	cmd := strings.TrimSpace(command) + " "

	commandMatcher, err := regexp.Compile("(.+?) ")
	if err != nil {
		return Command{}, err
	}

	args := commandMatcher.FindAllString(cmd, MaxParametersCount)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	if len(args) == 0 {
		return Command{}, errors.New("command for execution not found")
	}

	return Command{
		Cmd:        args[0],
		Parameters: args[1:],
	}, nil
}
