package server_session

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"net"
)

// ServerSession basic storage for server sessions
type ServerSession struct {
	Conn        *net.TCPConn
	AccessToken string
	Uploads     []*session.File
	Downloads   []*session.File
	Options     options.Options
}

// Release closes connection between server and client
func (s ServerSession) Release() error {
	if s.Conn == nil {
		return nil
	}

	err := s.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

// GetConn receives connection
func (s ServerSession) GetConn() *net.TCPConn {
	return s.Conn
}

// SetConn updates connection
func (s *ServerSession) SetConn(conn *net.TCPConn) {
	s.Conn = conn
}

// GetAccessToken receives client access token
func (s ServerSession) GetAccessToken() string {
	return s.AccessToken
}

// SetAccessToken updates client access token
func (s *ServerSession) SetAccessToken(token string) {
	s.AccessToken = token
}

// RegisterUpload initialize a File structure, append it to Uploads storage
func (s *ServerSession) RegisterUpload(filename string, filepath string) (*session.File, error) {
	file, err := session.CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
}

// RegisterDownload initialize a File structure, append it to Downloads storage
func (s ServerSession) RegisterDownload(filename string, filepath string) (*session.File, error) {
	file, err := session.CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
}

// FindUpload finds filename in Uploads slice
func (s ServerSession) FindUpload(filename string) *session.File {
	for _, file := range s.Uploads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

// FindDownload finds filename in Downloads slice
func (s ServerSession) FindDownload(filename string) *session.File {
	for _, file := range s.Downloads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

func (s ServerSession) GetOptions() options.Options {
	return s.Options
}

func (s ServerSession) ReceiveUnfinishedUploads() []string {
	uploads := make([]string, 0)
	for _, file := range s.Uploads {
		if !file.Completed() {
			uploads = append(uploads, file.Filename)
		}
	}
	return uploads
}

func (s ServerSession) ReceiveUnfinishedDownloads() []string {
	downloads := make([]string, 0)

	return downloads
}

// UploadStatus returns upload status in %
func (s ServerSession) UploadStatus() float64 {
	summary := float64(0)
	i := 0

	for _, file := range s.Uploads {
		if file.Completed() {
			continue
		}
		summary += float64(file.Size) / float64(file.Transferred)
		i++
	}

	if i == 0 {
		return 0
	}

	return summary / float64(i)
}

// DownloadStatus returns download status in %
func (s ServerSession) DownloadStatus() float64 {
	summary := float64(0)

	for _, file := range s.Downloads {
		summary += float64(file.Size) / float64(file.Transferred)
	}

	if len(s.Downloads) == 0 {
		return 0
	}

	return summary / float64(len(s.Downloads))
}

// CreateServerSession creates Session from connection and accessToken
func CreateServerSession(conn *net.TCPConn, options options.Options, accessToken string) session.Session {
	return &ServerSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*session.File{},
		Downloads:   []*session.File{},
		Options:     options,
	}
}
