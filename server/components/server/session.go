package server

import "net"

type SessionStorage interface {
	Find(host string) net.Conn
	Register(host string, conn net.Conn)
	Deregister(host string) error
}
