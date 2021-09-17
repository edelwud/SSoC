package executor

import (
	"errors"
	"server/components/command"
	"server/components/session"
)

type BasicExecService struct {
	serviceSet map[string]Executor
}

func (s BasicExecService) Process(session session.Session, cmd command.Command) error {
	executor := s.serviceSet[cmd.Cmd]
	if executor == nil {
		return errors.New("unrecognized command")
	}

	if !executor.CanAccess(session.GetAccessToken()) {
		return errors.New("cannot access to command: " + cmd.Cmd)
	}

	return executor.Process(session, cmd.Parameters...)
}

func RegisterBasicExecutorService(ctx session.SessionStorage) ExecService {
	execService := &BasicExecService{}
	execService.serviceSet = map[string]Executor{}
	execService.serviceSet["TIME"] = createTimeExecutor(ctx)
	execService.serviceSet["ECHO"] = createEchoExecutor(ctx)
	execService.serviceSet["CLOSE"] = createCloseExecutor(ctx)
	execService.serviceSet["TOKEN"] = createTokenExecutor(ctx)
	return execService
}
