package executor

import (
	"server/components/session"
)

type CloseExecutor struct {
	ctx session.SessionStorage
}

func (e CloseExecutor) Process(remoteAddr string, params ...string) error {
	err := e.ctx.Deregister(remoteAddr)
	if err != nil {
		return err
	}

	return nil
}

func createCloseExecutor(ctx session.SessionStorage) Executor {
	return &CloseExecutor{ctx: ctx}
}
