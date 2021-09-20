package session

import "net"

type Session interface {
	Release() error
	GetConn() *net.TCPConn
	SetConn(conn *net.TCPConn)
	GetAccessToken() string
	SetAccessToken(token string)
	RegisterUpload() *File
	RegisterDownload() *File
	FindUpload(filename string) *File
	FindDownload(filename string) *File
}
