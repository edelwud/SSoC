package session

import (
	"net"
)

type BasicSession struct {
	Conn        net.Conn
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

func (s BasicSession) GetConn() net.Conn {
	return s.Conn
}

func (s *BasicSession) SetConn(conn net.Conn) {
	s.Conn = conn
}

func (s BasicSession) GetAccessToken() string {
	return s.AccessToken
}

func (s *BasicSession) SetAccessToken(token string) {
	s.AccessToken = token
}

func (s *BasicSession) RegisterUpload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
}

func (s BasicSession) RegisterDownload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Uploads = append(s.Uploads, file)
	return file, nil
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

func CreateBasicSession(conn net.Conn, accessToken string) Session {
	return &BasicSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}
