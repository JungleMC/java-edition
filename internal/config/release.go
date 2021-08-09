//+build !dev

package config

type Config struct {
	DebugMode bool `env:"DEBUG" envDefault:"false"`
	Verbose   bool `env:"VERBOSE" envDefault:"false"`

	// Redis
	RedisHost string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	RedisPort int `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDatabase int`env:"REDIS_DB" envDefault:"0"`

	// Java Edition configuration
	ListenAddress        string `env:"LISTEN_ADDRESS" envDefault:"0.0.0.0"`
	ListenPort           int    `env:"LISTEN_PORT" envDefault:"25565"`
	OnlineMode           bool   `env:"ONLINE_MODE" envDefault:"true"`
	CompressionThreshold int    `env:"COMPRESSION_THRESHOLD" envDefault:"256"`
}
