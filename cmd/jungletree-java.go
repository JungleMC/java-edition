package main

import (
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net"
	"github.com/caarlos0/env"
)

func main() {
	config.Get = &config.Config{}
	if err := env.Parse(config.Get); err != nil {
		panic(err)
	}

	_, err := net.Bootstrap(config.Get.ListenAddress, config.Get.ListenPort, config.Get.OnlineMode)
	if err != nil {
		panic(err) // TODO: tidy error reporting?
	}
}
