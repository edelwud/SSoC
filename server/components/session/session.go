package session

import "net"

type Session interface {
	Release() error
	GetConn() net.Conn
	SetConn(conn net.Conn)
	GetAccessToken() string
	SetAccessToken(token string)
}
