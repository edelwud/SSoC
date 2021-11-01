package session

import (
	"net"
	"server/components/options"
)

// Session declares functionality for server sessions
type Session interface {
	Release() error
	GetOptions() options.Options
	GetConn() net.Conn
	SetConn(conn net.Conn)
	GetAccessToken() string
	SetAccessToken(token string)
	RegisterUpload(filename string, filepath string) (*File, error)
	RegisterDownload(filename string, filepath string) (*File, error)
	FindUpload(filename string) *File
	FindDownload(filename string) *File
}
