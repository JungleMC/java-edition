//+build dev

package config

type Config struct {
	DebugMode bool `env:"DEBUG" envDefault:"true"`
	Verbose   bool `env:"VERBOSE" envDefault:"true"`

	// Java Edition configuration
	ListenAddress        string `env:"LISTEN_ADDRESS" envDefault:"0.0.0.0"`
	ListenPort           int    `env:"LISTEN_PORT" envDefault:"25565"`
	OnlineMode           bool   `env:"ONLINE_MODE" envDefault:"true"`
	CompressionThreshold int    `env:"COMPRESSION_THRESHOLD" envDefault:"256"`
}
