package client_session

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
	"net"
)

// ClientSession basic storage for client sessions
type ClientSession struct {
	Conn        net.Conn
	AccessToken string
	Uploads     []*session.File
	Downloads   []*session.File
	Options     options.Options
}

// Release closes connection between server and client
func (s ClientSession) Release() error {
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
func (s ClientSession) GetConn() net.Conn {
	return s.Conn
}

// SetConn updates connection
func (s *ClientSession) SetConn(conn net.Conn) {
	s.Conn = conn
}

// GetAccessToken receives client access token
func (s ClientSession) GetAccessToken() string {
	return s.AccessToken
}

// SetAccessToken updates client access token
func (s *ClientSession) SetAccessToken(token string) {
	s.AccessToken = token
}

// FreeUpload finds upload and erases them
func (s ClientSession) FreeUpload(filename string) {
	for i, file := range s.Uploads {
		if file.Filename == filename {
			s.Uploads = append(s.Uploads[:i], s.Uploads[i+1:]...)
		}
	}
}

// RegisterUpload initialize a File structure, append it to Uploads storage
func (s *ClientSession) RegisterUpload(filename string, filepath string) (*session.File, error) {
	s.FreeUpload(filename)

	upload, err := session.CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}

	s.Uploads = append(s.Uploads, upload)
	return upload, nil
}

// RegisterDownload initialize a File structure, append it to Downloads storage
func (s *ClientSession) RegisterDownload(filename string, filepath string) (*session.File, error) {
	file, err := session.CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Downloads = append(s.Downloads, file)
	return file, nil
}

// FindUpload finds filename in Uploads slice
func (s ClientSession) FindUpload(filename string) *session.File {
	for _, file := range s.Uploads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

// FindDownload finds filename in Downloads slice
func (s ClientSession) FindDownload(filename string) *session.File {
	for _, file := range s.Downloads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

// UploadStatus returns upload status in %
func (s ClientSession) UploadStatus() float64 {
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
func (s ClientSession) DownloadStatus() float64 {
	summary := float64(0)

	for _, file := range s.Downloads {
		summary += float64(file.Size) / float64(file.Transferred)
	}

	if len(s.Downloads) == 0 {
		return 0
	}

	return summary / float64(len(s.Downloads))
}

// GetOptions returns options for connection to server
func (s ClientSession) GetOptions() options.Options {
	return s.Options
}

func (s ClientSession) ReceiveUnfinishedUploads() []string {
	uploads := make([]string, 0)
	for _, file := range s.Uploads {
		if !file.Completed() {
			uploads = append(uploads, file.Filename)
		}
	}
	return uploads
}

func (s ClientSession) ReceiveUnfinishedDownloads() []string {
	downloads := make([]string, 0)

	return downloads
}

// CreateClientSession creates Session from connection and accessToken
func CreateClientSession(conn net.Conn, options options.Options, accessToken string) session.Session {
	return &ClientSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*session.File{},
		Downloads:   []*session.File{},
		Options:     options,
	}
}
