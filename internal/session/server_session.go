package session

import (
	"SSoC/internal/options"
	"net"
)

// ServerSession basic storage for server sessions
type ServerSession struct {
	Conn        net.Conn
	Addr        net.Addr
	AccessToken string
	Uploads     []*File
	Downloads   []*File
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
func (s ServerSession) GetConn() net.Conn {
	return s.Conn
}

// SetConn updates connection
func (s *ServerSession) SetConn(conn net.Conn) {
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
func (s *ServerSession) RegisterUpload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
}

// RegisterDownload initialize a File structure, append it to Downloads storage
func (s ServerSession) RegisterDownload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
}

// FindUpload finds filename in Uploads slice
func (s ServerSession) FindUpload(filename string) *File {
	for _, file := range s.Uploads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

// FindDownload finds filename in Downloads slice
func (s ServerSession) FindDownload(filename string) *File {
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

func (s ServerSession) GetAddress() net.Addr {
	return s.Addr
}

// CreateServerSession creates Session from connection and accessToken
func CreateServerSession(conn net.Conn, options options.Options, accessToken string, addr net.Addr) Session {
	return &ServerSession{
		Addr:        addr,
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
		Options:     options,
	}
}
