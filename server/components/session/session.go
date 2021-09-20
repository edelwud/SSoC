package session

import "net"

type Session interface {
	Release() error
	GetConn() net.Conn
	SetConn(conn net.Conn)
	GetAccessToken() string
	SetAccessToken(token string)
	RegisterUpload() *File
	RegisterDownload() *File
	FindUpload(filename string) *File
	FindDownload(filename string) *File
}
