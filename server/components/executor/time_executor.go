package executor

import (
	"server/components/session"
	"time"
)

type TimeExecutor struct {
	ctx session.SessionStorage
}

func (e TimeExecutor) Process(remoteAddr string, _ ...string) error {
	conn, err := e.ctx.Find(remoteAddr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(time.Now().String() + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func createTimeExecutor(ctx session.SessionStorage) Executor {
	return &TimeExecutor{ctx: ctx}
}
