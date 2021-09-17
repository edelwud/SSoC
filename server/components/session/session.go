package session

import "net"

type SessionStorage interface {
	Find(host string) (net.Conn, error)
	Register(host string, conn net.Conn)
	Deregister(host string) error
}
