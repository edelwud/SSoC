package server

type Server interface {
	Open() error
	Close() error
}

type Options struct {
	Host string
	Port string
}
