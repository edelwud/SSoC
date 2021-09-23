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

func (s BasicSession) FreeUpload(filename string) {
	for i, file := range s.Uploads {
		if file.Filename == filename {
			s.Uploads = append(s.Uploads[:i], s.Uploads[i+1:]...)
		}
	}
}

func (s *BasicSession) RegisterUpload(filename string, filepath string) (*File, error) {
	s.FreeUpload(filename)

	upload, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}

	s.Uploads = append(s.Uploads, upload)
	return upload, nil
}

func (s *BasicSession) RegisterDownload(filename string, filepath string) (*File, error) {
	file, err := CreateFile(filename, filepath)
	if err != nil {
		return nil, err
	}
	s.Downloads = append(s.Downloads, file)
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

func (s BasicSession) UploadStatus() float64 {
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

func (s BasicSession) DownloadStatus() float64 {
	summary := float64(0)

	for _, file := range s.Downloads {
		summary += float64(file.Size) / float64(file.Transferred)
	}

	if len(s.Downloads) == 0 {
		return 0
	}

	return summary / float64(len(s.Downloads))
}

func CreateBasicSession(conn *net.TCPConn, accessToken string) Session {
	return &BasicSession{
		Conn:        conn,
		AccessToken: accessToken,
		Uploads:     []*File{},
		Downloads:   []*File{},
	}
}
