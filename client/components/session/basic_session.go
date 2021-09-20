package session

import (
	"net"
)

type BasicSession struct {
	Conn        *net.TCPConn
	AccessToken string
	Uploads     []*File
	Downloads   []*File
}

func (s BasicSession) Release() error {
	if s.Conn == nil {
		return nil
	}

	err := s.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s BasicSession) GetConn() *net.TCPConn {
	return s.Conn
}

func (s *BasicSession) SetConn(conn *net.TCPConn) {
	s.Conn = conn
}

func (s BasicSession) GetAccessToken() string {
	return s.AccessToken
}

func (s *BasicSession) SetAccessToken(token string) {
	s.AccessToken = token
}

func (s *BasicSession) RegisterUpload() *File {
	file := &File{}
	s.Uploads = append(s.Uploads, file)
	return file
}

func (s BasicSession) RegisterDownload() *File {
	file := &File{}
	s.Downloads = append(s.Downloads, file)
	return file
}

func (s BasicSession) FindUpload(filename string) *File {
	for _, file := range s.Uploads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

func (s BasicSession) FindDownload(filename string) *File {
	for _, file := range s.Downloads {
		if file.Filename == filename {
			return file
		}
	}
	return nil
}

func CreateBasicSession(conn *net.TCPConn, accessToken string) Session {
	return &BasicSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}
