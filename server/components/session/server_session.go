package session

import (
	"net"
)

// ServerSession basic storage for server sessions
type ServerSession struct {
	Conn        net.Conn
	AccessToken string
	Uploads     []*File
	Downloads   []*File
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

// CreateServerSession creates Session from connection and accessToken
func CreateServerSession(conn net.Conn, accessToken string) Session {
	return &ServerSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}
