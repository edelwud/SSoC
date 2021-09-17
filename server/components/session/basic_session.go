package session

import "net"

type BasicSession struct {
	Conn        net.Conn
	AccessToken string
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

func CreateBasicSession(conn net.Conn, accessToken string) Session {
	return &BasicSession{Conn: conn, AccessToken: accessToken}
}
