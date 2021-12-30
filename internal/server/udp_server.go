package server

import (
	"SSoC/internal/command"
	"SSoC/internal/executor"
	"SSoC/internal/options"
	"SSoC/internal/session"
	"errors"
	"fmt"
	"io"
	"net"
)

// UDPServer implementation of Server interfaces based on UDP protocol
type UDPServer struct {
	Options     options.Options
	Conn        *net.UDPConn
	Context     session.Storage
	ExecService executor.Service
	Clients     map[string]session.Session
	UDPWriter   *UDPWriter
}

type UDPWriter struct {
	Addr net.Addr
	Conn *net.UDPConn
}

const CommandBufferSize = 1024 * 1024

// Run resolves server options from Options
// creates net.Listener with UDPv4 background
// executes AcceptLoop
func (s *UDPServer) Run() error {
	addr, err := net.ResolveUDPAddr("udp", s.Options.Host+":"+s.Options.Port)
	if err != nil {
		return err
	}

	s.Conn, err = net.ListenUDP("udp", addr)
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

func (s UDPServer) Auth(writer io.Writer, addr net.Addr, c string) error {
	cmd, err := command.ParseCommand(c)
	if err != nil {
		return err
	}

	if cmd.Cmd != "TOKEN" {
		return errors.New("cannot authenticate new client")
	}

	sess := session.CreateServerSession(s.Conn, s.Options, cmd.Parameters[0], addr)

	err = s.AddClient(addr.String(), sess)
	if err != nil {
		return err
	}

	err = s.ExecService.Process(writer, sess, cmd)
	if err != nil {
		return err
	}

	return nil
}

// AcceptLoop accepts client connection, sets keep alive and keep alive period options from Options, handles connection
func (s *UDPServer) AcceptLoop() error {
	for {
		buf := make([]byte, CommandBufferSize)
		n, addr, err := s.Conn.ReadFromUDP(buf)
		if err != nil {
			serverLogger.Warnf("client disconnected")
			err := s.Run()
			if err != nil {
				return err
			}
			return nil
		}

		go s.HandleClient(addr, buf[:n])
	}
}

func (s UDPServer) HandleClient(addr net.Addr, buf []byte) {
	fmt.Println(addr.String())

	sess, err := s.FindClient(addr.String())
	writer := CreateUDPWriter(addr, s.Conn)
	if err != nil {
		err := s.Auth(writer, addr, string(buf))
		if err != nil {
			serverLogger.Warnf("client disconnected")
			return
		}
		return
	}

	cmd, err := command.ParseCommand(string(buf))
	if err != nil {
		serverLogger.Warnf("sent undefined command")
		return
	}

	err = s.ExecService.Process(writer, sess, cmd)
	if err != nil {
		serverLogger.Warnf("command execution error: %s", err)
		return
	}

	serverLogger.Infof("command %s successfully executed", cmd.Cmd)
}

func (s UDPServer) FindClient(address string) (session.Session, error) {
	if _, ok := s.Clients[address]; !ok {
		return nil, errors.New("client not found")
	}

	return s.Clients[address], nil
}

func (s UDPServer) AddClient(address string, session session.Session) error {
	if _, ok := s.Clients[address]; ok {
		return errors.New("client already exists")
	}

	s.Clients[address] = session
	return nil
}

// Close closes net.Listener
func (s *UDPServer) Close() error {
	if s.Conn == nil {
		return nil
	}

	err := s.Conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (w UDPWriter) Write(p []byte) (n int, err error) {
	n, err = w.Conn.WriteTo(p, w.Addr)
	if err != nil {
		return 0, err
	}

	return n, err
}

func CreateUDPWriter(addr net.Addr, conn *net.UDPConn) *UDPWriter {
	return &UDPWriter{
		Addr: addr,
		Conn: conn,
	}
}

// CreateUDPServer creates UDPServer with initialized session.ServerStorage and executor.ServerExecutorService
func CreateUDPServer(options options.Options) Server {
	ctx := session.CreateServerSessionStorage()
	executorService := executor.RegisterServerExecutorService(ctx)
	return &UDPServer{
		Options:     options,
		Context:     ctx,
		ExecService: executorService,
		Clients:     make(map[string]session.Session),
	}
}
