package datachannel

import (
	"SSoC/internal/options"
	"math/rand"
	"strconv"
)

func generatePort() string {
	return strconv.Itoa(int(rand.Float32()*1000) + 8000)
}

var dataChannelServer = map[string]func(options.Options) Datachannel{
	"udp": func(ops options.Options) Datachannel {
		return NewUDPDatachannelServer(ops)
	},
	"tcp": func(ops options.Options) Datachannel {
		return NewTCPDatachannelServer(ops)
	},
}

func NewServer(protocol string, ops options.Options) Datachannel {
	if channel, ok := dataChannelServer[protocol]; !ok {
		return nil
	} else {
		return channel(ops)
	}
}
