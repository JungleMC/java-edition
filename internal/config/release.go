//+build !dev

package config

var Get *Config

type Config struct {
	DebugMode bool `env:"DEBUG" envDefault:"false"`
	Verbose bool `env:"VERBOSE" envDefault:"false"`
	// Java Edition configuration
	ListenAddress string `env:"LISTEN_ADDRESS" envDefault:"0.0.0.0"`
	ListenPort int `env:"LISTEN_PORT" envDefault:"25565"`
	OnlineMode bool `env:"ONLINE_MODE" envDefault:"true"`
	CompressionThreshold int `env:"COMPRESSION_THRESHOLD" envDefault:"256"`

	// RPC services
	StatusHost string `env:"STATUS_HOST" envDefault:"127.0.0.1"`
	StatusPort int `env:"STATUS_PORT" envDefault:"50050"`
}
