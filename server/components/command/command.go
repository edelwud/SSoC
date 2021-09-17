package command

import (
	"errors"
	"regexp"
	"strings"
)

type Command struct {
	Cmd        string
	Parameters []string
}

const MaxParametersCount = 100

func ParseCommand(command string) (Command, error) {
	command = strings.TrimSpace(command) + " "

	commandMatcher, err := regexp.Compile("(.+?) ")
	if err != nil {
		return Command{}, err
	}

	args := commandMatcher.FindAllString(command, MaxParametersCount)
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
