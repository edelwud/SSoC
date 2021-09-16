package server

import (
	"bufio"
	"errors"
	"net"
	"server/components/command"
	"strings"
	"time"
)

type TcpServer struct {
	Host     string
	Port     string
	Listener *net.TCPListener
	Session  SessionStorage
}

func (s *TcpServer) Run() error {
	addr, err := net.ResolveTCPAddr("tcp", s.Host+":"+s.Port)
	if err != nil {
		return err
	}

	s.Listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	serverLogger.Infof("server started on port %s", s.Port)

	err = s.AcceptLoop()
	if err != nil {
		return err
	}

	return nil
}

func (s *TcpServer) AcceptLoop() error {
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			return err
		}

		serverLogger.Infof("client %s connected", conn.RemoteAddr())

		err = conn.SetKeepAlive(true)
		if err != nil {
			return err
		}

		err = conn.SetKeepAlivePeriod(time.Minute)
		if err != nil {
			return err
		}

		if err != nil {
			return err
		}
		go s.HandleConnection(conn)
	}
}

func (s *TcpServer) HandleConnection(conn net.Conn) {
	for {
		userCommand, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			serverLogger.Infof("client %s disconnected, reason: %s", conn.RemoteAddr(), err)
			return
		}

		cmd, err := command.ParseCommand(userCommand)
		if err != nil {
			serverLogger.Warnf("client %s sent undefined command", conn.RemoteAddr())
			continue
		}

		err = s.ExecuteCommand(conn, cmd)
		if err != nil {
			if err.Error() == "close connection interrupt" {
				serverLogger.Infof("client %s disconnected, reason: %s", conn.RemoteAddr(), err)
			} else {
				serverLogger.Warnf("client %s command execution error: %s", conn.RemoteAddr(), err)
			}
			return
		}
	}
}

func (s *TcpServer) ExecuteCommand(conn net.Conn, cmd *command.Command) error {
	switch cmd.Execute {
	case command.EchoExec:
		_, err := conn.Write([]byte(strings.Join(cmd.Parameters, "") + "\n"))
		if err != nil {
			return err
		}
	case command.TimeExec:
		t := time.Now()
		now := t.Format(time.RFC3339) + "\n"
		_, err := conn.Write([]byte(now))
		if err != nil {
			return err
		}
	case command.CloseConnectionExec:
		err := conn.Close()
		if err != nil {
			return err
		}
		return errors.New("close connection interrupt")
	}
	return nil
}

func (s *TcpServer) Close() error {
	if s.Listener == nil {
		return nil
	}

	err := s.Listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func CreateTcpServer(options Options) Server {
	return &TcpServer{Host: options.Host, Port: options.Port}
}
