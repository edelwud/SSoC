package executor

import (
	"server/components/session"
)

type Executor interface {
	Process(remoteAddr string, params ...string) error
}

func RegisterExecutorService(ctx session.SessionStorage) ExecService {
	execService := &BasicExecService{}
	execService.serviceSet = map[string]Executor{}
	execService.serviceSet["TIME"] = createTimeExecutor(ctx)
	execService.serviceSet["ECHO"] = createEchoExecutor(ctx)
	execService.serviceSet["CLOSE"] = createCloseExecutor(ctx)
	return execService
}
