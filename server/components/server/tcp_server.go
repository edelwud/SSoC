package server

import (
	"bufio"
	"net"
	"server/components/command"
	"server/components/executor"
	"server/components/session"
	"time"
)

type TcpServer struct {
	Options     Options
	Listener    *net.TCPListener
	Context     session.SessionStorage
	ExecService executor.ExecService
}

func (s *TcpServer) Run() error {
	addr, err := net.ResolveTCPAddr("tcp", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}

	s.Listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	serverLogger.Infof("server started on port %s", s.Options.Port)

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

		connectionLogger := serverLogger.WithField("client", conn.RemoteAddr())
		connectionLogger.Info("connected")

		err = conn.SetKeepAlive(s.Options.KeepAlive)
		if err != nil {
			return err
		}

		err = conn.SetKeepAlivePeriod(time.Second * time.Duration(s.Options.KeepAlivePeriod))
		if err != nil {
			return err
		}

		go s.HandleConnection(conn)
	}
}

func (s *TcpServer) HandleConnection(conn net.Conn) {
	connectionLogger := serverLogger.WithField("client", conn.RemoteAddr())

	accessToken := ""
	currentSession := session.CreateBasicSession(conn, accessToken)

	for {
		userCommand, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			connectionLogger.Infof("disconnected, reason: %s", err)
			return
		}

		cmd, err := command.ParseCommand(userCommand)
		if err != nil {
			connectionLogger.Warnf("sent undefined command")
			continue
		}

		err = s.ExecService.Process(currentSession, cmd)
		if err != nil {
			connectionLogger.Warnf("command execution error: %s", err)
			continue
		}

		connectionLogger.Infof("command %s successfully executed", cmd.Cmd)
	}
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
	ctx := session.CreateBasicSessionStorage()
	executorService := executor.RegisterBasicExecutorService(ctx)
	return &TcpServer{Options: options, Context: ctx, ExecService: executorService}
}
