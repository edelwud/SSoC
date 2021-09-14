package command

import (
	"errors"
	"regexp"
	"strings"
)

type ExecCommand int64

type Command struct {
	Execute    ExecCommand
	Parameters []string
}

const (
	EchoExec ExecCommand = iota
	UploadExec
	DownloadExec
	CloseConnectionExec
	UndefinedCommand
)

const MaxParametersCount = 100

func TransformExec(exec string) (ExecCommand, error) {
	switch exec {
	case "ECHO":
		return EchoExec, nil
	case "UPLOAD":
		return UploadExec, nil
	case "DOWNLOAD":
		return DownloadExec, nil
	case "CLOSE":
		return CloseConnectionExec, nil
	}
	return UndefinedCommand, errors.New("cannot recognize execution command")
}

func ParseCommand(command string) (*Command, error) {
	command = strings.TrimSpace(command) + " "

	commandMatcher, err := regexp.Compile("(.+?) ")
	if err != nil {
		return nil, err
	}

	args := commandMatcher.FindAllString(command, MaxParametersCount)
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}

	execute, err := TransformExec(args[0])
	if err != nil {
		return nil, err
	}

	return &Command{
		Execute:    execute,
		Parameters: args[1:],
	}, nil
}
