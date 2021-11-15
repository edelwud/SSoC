package datachannel

import "SSoC/internal/session"

type Datachannel interface {
	Connect() error
	Close() error
	Upload(file *session.File) error
	Download(file *session.File) error
}
