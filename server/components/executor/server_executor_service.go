package executor

import (
	"errors"
	"server/components/command"
	"server/components/session"
)

// ServerExecutorService stores command executors in ServiceSet
type ServerExecutorService struct {
	ServiceSet map[string]Executor
}

// Process finds responsible command in ServiceSet, checks client access and executes them
func (s ServerExecutorService) Process(session session.Session, cmd command.Command) error {
	executor := s.ServiceSet[cmd.Cmd]
	if executor == nil {
		return errors.New("unrecognized command")
	}

	if !executor.CanAccess(session.GetAccessToken()) {
		return errors.New("cannot access to command: " + cmd.Cmd)
	}

	return executor.Process(session, cmd.Parameters...)
}

// RegisterServerExecutorService register ServerExecutorService.ServiceSet with all responsible commands
func RegisterServerExecutorService(ctx session.Storage) Service {
	execService := &ServerExecutorService{}
	execService.ServiceSet = map[string]Executor{}
	execService.ServiceSet["TIME"] = createTimeExecutor(ctx)
	execService.ServiceSet["ECHO"] = createEchoExecutor(ctx)
	execService.ServiceSet["CLOSE"] = createCloseExecutor(ctx)
	execService.ServiceSet["TOKEN"] = createTokenExecutor(ctx)
	execService.ServiceSet["UPLOAD"] = createUploadExecutor(ctx)
	execService.ServiceSet["DOWNLOAD"] = createDownloadExecutor(ctx)
	execService.ServiceSet["REQUEST_UPLOAD"] = createRequestUploadExecutor(ctx)
	execService.ServiceSet["REQUEST_DOWNLOAD"] = createRequestDownloadExecutor(ctx)
	return execService
}
