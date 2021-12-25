package options

// Options declares generalized structure of server parameters
type Options struct {
	Protocol        string `yaml:"protocol"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	KeepAlive       bool   `yaml:"keepAlive"`
	KeepAlivePeriod int    `yaml:"keepAlivePeriod"`
}
