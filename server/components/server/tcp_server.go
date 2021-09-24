package server

import (
	"bufio"
	"net"
	"server/components/command"
	"server/components/executor"
	"server/components/session"
	"time"
)

// TcpServer implementation of Server interfaces based on TCP protocol
type TcpServer struct {
	Options     Options
	Listener    *net.TCPListener
	Context     session.Storage
	ExecService executor.Service
}

// Run resolves server options from Options
// creates net.Listener with TCPv4 background
// executes AcceptLoop
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

// AcceptLoop accepts client connection, sets keep alive and keep alive period options from Options, handles connection
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

// HandleConnection creates server session for each connection, reads client command and executes it
func (s *TcpServer) HandleConnection(conn net.Conn) {
	connectionLogger := serverLogger.WithField("client", conn.RemoteAddr())

	accessToken := ""
	currentSession := session.CreateServerSession(conn, accessToken)

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

// Close closes net.Listener
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

// CreateTcpServer creates TcpServer with initialized session.ServerStorage and executor.ServerExecutorService
func CreateTcpServer(options Options) Server {
	ctx := session.CreateServerSessionStorage()
	executorService := executor.RegisterServerExecutorService(ctx)
	return &TcpServer{
		Options:     options,
		Context:     ctx,
		ExecService: executorService,
	}
}
