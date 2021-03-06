package server

import (
	"SSoC/internal/command"
	"SSoC/internal/executor"
	"SSoC/internal/options"
	"SSoC/internal/session"
	"bufio"
	"log"
	"net"
	"syscall"
	"time"
)

// TCPServer implementation of Server interfaces based on TCP protocol
type TCPServer struct {
	Options     options.Options
	Listener    *net.TCPListener
	Context     session.Storage
	ExecService executor.Service
}

var epoller *epoll

func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}

// Run resolves server options from Options
// creates net.Listener with TCPv4 background
// executes AcceptLoop
func (s *TCPServer) Run() error {
	setLimit()

	addr, err := net.ResolveTCPAddr("tcp", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}

	s.Listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	epoller, err = MkEpoll()
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
func (s *TCPServer) AcceptLoop() error {
	for {
		conn, err := s.Listener.AcceptTCP()
		if err != nil {
			return err
		}

		err = conn.SetNoDelay(true)
		if err != nil {
			return err
		}

		err = epoller.Add(conn)
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
func (s *TCPServer) HandleConnection(conn *net.TCPConn) {
	connectionLogger := serverLogger.WithField("client", conn.RemoteAddr())

	accessToken := ""
	currentSession := session.CreateServerSession(conn, s.Options, accessToken)

	for {
		connections, err := epoller.Wait()
		if err != nil {
			connectionLogger.Infof("disconnected, reason: %s", err)
		}

		for _, conn := range connections {
			userCommand, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				connectionLogger.Infof("disconnected, reason: %s", err)
				return
			}

			cmd, err := command.ParseCommand(userCommand)
			if err != nil {
				connectionLogger.Warnf("sent undefined command")
				continue
			}

			err = s.ExecService.Process(conn, currentSession, cmd)
			if err != nil {
				connectionLogger.Warnf("command execution error: %s", err)
				continue
			}

			connectionLogger.Infof("command %s successfully executed", cmd.Cmd)
		}
	}
}

// Close closes net.Listener
func (s *TCPServer) Close() error {
	if s.Listener == nil {
		return nil
	}

	err := s.Listener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s TCPServer) Write(payload, token string) error {
	sess, err := s.Context.Find(token)
	if err != nil {
		return err
	}

	_, err = sess.GetConn().Write([]byte(payload + "\n"))
	if err != nil {
		return err
	}

	return nil
}

// CreateTCPServer creates TCPServer with initialized session.ServerStorage and executor.ServerExecutorService
func CreateTCPServer(options options.Options) Server {
	ctx := session.CreateServerSessionStorage()
	executorService := executor.RegisterServerExecutorService(ctx)
	return &TCPServer{
		Options:     options,
		Context:     ctx,
		ExecService: executorService,
	}
}
