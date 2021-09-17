package session

type SessionStorage interface {
	Find(accessToken string) (Session, error)
	Register(ctx Session)
	Deregister(accessToken string) error
}
