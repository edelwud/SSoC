package session

import (
	"net"
)

// ClientSession basic storage for client sessions
type ClientSession struct {
	Conn        *net.TCPConn
	AccessToken string
	Uploads     []*File
	Downloads   []*File
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
func (s ClientSession) GetConn() *net.TCPConn {
	return s.Conn
}

// SetConn updates connection
func (s *ClientSession) SetConn(conn *net.TCPConn) {
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
func (s *ClientSession) RegisterUpload(filename string, filepath string) (*File, error) {
	s.FreeUpload(filename)

	upload, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}

	s.Uploads = append(s.Uploads, upload)
	return upload, nil
}

// RegisterDownload initialize a File structure, append it to Downloads storage
func (s *ClientSession) RegisterDownload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Downloads = append(s.Downloads, file)
	return file, nil
}

// FindUpload finds filename in Uploads slice
func (s ClientSession) FindUpload(filename string) *File {
	for _, file := range s.Uploads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

// FindDownload finds filename in Downloads slice
func (s ClientSession) FindDownload(filename string) *File {
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

// CreateClientSession creates Session from connection and accessToken
func CreateClientSession(conn *net.TCPConn, accessToken string) Session {
	return &ClientSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}
