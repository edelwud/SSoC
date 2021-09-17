package executor

import (
	"server/components/session"
	"strings"
)

type EchoExecutor struct {
	ctx session.SessionStorage
}

func (e EchoExecutor) Process(remoteAddr string, params ...string) error {
	conn, err := e.ctx.Find(remoteAddr)
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte(strings.Join(params, " ") + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func createEchoExecutor(ctx session.SessionStorage) Executor {
	return &EchoExecutor{ctx: ctx}
}
