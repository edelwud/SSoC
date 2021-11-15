package datachannel

import "SSoC/internal/session"

type Datachannel interface {
	Listen() error
	Close() error
	Accept() error
	GetPort() string
	Upload(file *session.File) error
	Download(file *session.File) error
}
