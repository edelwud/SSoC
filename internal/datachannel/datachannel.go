package datachannel

import (
	"SSoC/internal/options"
	"SSoC/internal/session"
)

type Datachannel interface {
	Connect() error
	Listen() error
	Close() error
	Accept() error
	GetPort() string
	Upload(file *session.File) error
	Download(file *session.File) error
}

var dataChannels = map[string]func(string, string, options.Options) Datachannel{
	"server": func(protocol, port string, ops options.Options) Datachannel {
		return NewServer(protocol, ops)
	},
	"client": func(protocol, port string, ops options.Options) Datachannel {
		return NewClient(protocol, port, ops)
	},
}

func New(side, protocol, port string, ops options.Options) Datachannel {
	if channel, ok := dataChannels[side]; !ok {
		return nil
	} else {
		return channel(protocol, port, ops)
	}
}
