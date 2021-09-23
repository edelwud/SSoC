package session

import "net"

type Session interface {
	Release() error
	GetConn() *net.TCPConn
	SetConn(conn *net.TCPConn)
	GetAccessToken() string
	SetAccessToken(token string)
	RegisterUpload(filename string, filepath string) (*File, error)
	RegisterDownload(filename string, filepath string) (*File, error)
	FindUpload(filename string) *File
	FindDownload(filename string) *File
}
