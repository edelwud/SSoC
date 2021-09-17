package session

import (
	"errors"
	"net"
)

type BasicSessionStorage struct {
	clients map[string]net.Conn
}

func (s BasicSessionStorage) Find(host string) (net.Conn, error) {
	conn := s.clients[host]
	if conn == nil {
		return nil, errors.New("client not found")
	}
	return conn, nil
}

func (s *BasicSessionStorage) Register(host string, conn net.Conn) {
	s.clients[host] = conn
}

func (s *BasicSessionStorage) Deregister(host string) error {
	if s.clients[host] == nil {
		return nil
	}

	err := s.clients[host].Close()
	if err != nil {
		return err
	}

	s.clients[host] = nil

	return nil
}

func CreateBasicSessionStorage() SessionStorage {
	session := &BasicSessionStorage{}
	session.clients = map[string]net.Conn{}
	return session
}
