package executor

import (
	"server/components/session"
)

type TimeExecutor struct {
	ctx session.SessionStorage
}

func (e TimeExecutor) Process(remoteAddr string, params ...string) error {
	return nil
}

func createTimeExecutor(ctx session.SessionStorage) Executor {
	return &TimeExecutor{ctx: ctx}
}
