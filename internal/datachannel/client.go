package datachannel

import (
	"SSoC/internal/options"
)

var dataChannelsClient = map[string]func(string, options.Options) Datachannel{
	"udp": func(port string, ops options.Options) Datachannel {
		return NewUDPDatachannelClient(port, ops)
	},
	"tcp": func(port string, ops options.Options) Datachannel {
		return NewTCPDatachannelClient(port, ops)
	},
}

func NewClient(protocol string, port string, ops options.Options) Datachannel {
	if channel, ok := dataChannelsClient[protocol]; !ok {
		return nil
	} else {
		return channel(port, ops)
	}
}
