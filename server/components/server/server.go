package server

type Server interface {
	Run() error
	Close() error
}

type Options struct {
	Host string
	Port string
}
