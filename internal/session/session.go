package session

import (
	"SSoC/internal/options"
	"net"
)

// Session declares functionality for server sessions
type Session interface {
	Release() error
	ReceiveUnfinishedUploads() []string
	ReceiveUnfinishedDownloads() []string

	UploadStatus() float64
	DownloadStatus() float64

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
