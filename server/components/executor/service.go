package executor

import (
	"errors"
	"server/components/command"
)

type ExecService interface {
	Process(remoteAddr string, cmd command.Command) error
}

type BasicExecService struct {
	serviceSet map[string]Executor
}

func (s BasicExecService) Process(remoteAddr string, cmd command.Command) error {
	executor := s.serviceSet[cmd.Cmd]
	if executor == nil {
		return errors.New("unrecognized command")
	}
	return executor.Process(remoteAddr, cmd.Parameters...)
}
