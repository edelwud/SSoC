package session

import (
	"errors"
)

type BasicSessionStorage struct {
	clients map[string]Session
}

func (s BasicSessionStorage) Find(accessToken string) (Session, error) {
	conn := s.clients[accessToken]
	if conn == nil {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

func (s *BasicSessionStorage) Register(ctx Session) {
	s.clients[ctx.GetAccessToken()] = ctx
}

func (s *BasicSessionStorage) Deregister(accessToken string) error {
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

func CreateBasicSessionStorage() SessionStorage {
	session := &BasicSessionStorage{}
	session.clients = map[string]Session{}
	return session
}
