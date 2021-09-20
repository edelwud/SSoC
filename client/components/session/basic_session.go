package session

import "net"

type File struct {
	Filename    string
	Transferred uint
	Size        uint
}

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

func (s BasicSession) RegisterUpload() *File {
	file := &File{}
	s.Uploads = append(s.Uploads, file)
	return file
}

func (s BasicSession) RegisterDownload() *File {
	file := &File{}
	s.Downloads = append(s.Downloads, file)
	return file
}

func CreateBasicSession(conn net.Conn, accessToken string) Session {
	return &BasicSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}