package session

import (
	"errors"
)

// ServerStorage stores clients sessions
type ServerStorage struct {
	clients map[string]Session
}

// Find searches accessToken in clients map
func (s ServerStorage) Find(accessToken string) (Session, error) {
	conn := s.clients[accessToken]
	if conn == nil {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

// Register saves Session in clients map where key is accessToken
func (s *ServerStorage) Register(ctx Session) {
	s.clients[ctx.GetAccessToken()] = ctx
}

// Deregister removes Session from clients map, executes Session.Release for closing connection
func (s *ServerStorage) Deregister(accessToken string) error {
	if s.clients[accessToken] == nil {
		return nil
	}

	err := s.clients[accessToken].Release()
	if err != nil {
		return err
	}

	s.clients[accessToken] = nil

	return nil
}

// CreateServerSessionStorage creates server session storage
func CreateServerSessionStorage() Storage {
	return &ServerStorage{
		clients: map[string]Session{},
	}
}
