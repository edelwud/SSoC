package executor

import (
	"server/components/session"
)

type Executor interface {
	Process(session session.Session, params ...string) error
	CanAccess(accessToken string) bool
}
