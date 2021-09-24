package session

// Storage declares functionality for accessing to Session
type Storage interface {
	Find(accessToken string) (Session, error)
	Register(ctx Session)
	Deregister(accessToken string) error
}
