package server

// Server declares generalized interface for server functionality
type Server interface {
	Run() error
	Close() error
}
