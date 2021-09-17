package executor

import (
	"server/components/session"
)

type EchoExecutor struct {
	ctx session.SessionStorage
}

func (e EchoExecutor) Process(remoteAddr string, params ...string) error {
	return nil
}

func createEchoExecutor(ctx session.SessionStorage) Executor {
	return &EchoExecutor{ctx: ctx}
}
