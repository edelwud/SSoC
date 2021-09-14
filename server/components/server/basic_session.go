package server

import "net"

type BasicSessionStorage struct {
	clients map[string]net.Conn
}

func (s BasicSessionStorage) Find(host string) net.Conn {
	return s.clients[host]
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

	return nil
}

func CreateBasicSessionStorage() SessionStorage {
	return &BasicSessionStorage{}
}
