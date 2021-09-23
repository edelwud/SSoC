package session

import "net"

type Session interface {
	Release() error
	GetConn() *net.TCPConn
	SetConn(conn *net.TCPConn)
	GetAccessToken() string
	SetAccessToken(token string)

	UploadStatus() float64
	FindUpload(filename string) *File
	RegisterUpload(filename string, filepath string) (*File, error)

	DownloadStatus() float64
	FindDownload(filename string) *File
	RegisterDownload(filename string, filepath string) (*File, error)
}
