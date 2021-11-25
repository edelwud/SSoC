package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
)

type Datachannel interface {
	Connect() error
	Close() error
	Upload(file *session.File) error
	Download(file *session.File) error
}

var dc = map[string]func(string, options.Options) Datachannel{
	"tcp": func(port string, ops options.Options) Datachannel {
		return NewTCPDatachannel(port, ops)
	},
	"udp": func(port string, ops options.Options) Datachannel {
		return NewUDPDatachannel(port, ops)
	},
}

func NewDatachannel(port string, ops options.Options) Datachannel {
	return dc[ops.Protocol](port, ops)
}
