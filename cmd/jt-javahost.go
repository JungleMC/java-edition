package main

import (
	"github.com/caarlos0/env"
	"github.com/junglemc/Service-JavaEditionHost/internal/config"
	"github.com/junglemc/Service-JavaEditionHost/internal/net"
	"github.com/junglemc/Service-JavaEditionHost/internal/rpc"
)

func main() {
	config.Get = &config.Config{}
	if err := env.Parse(config.Get); err != nil {
		panic(err)
	}

	rpc.StatusInit(config.Get.StatusHost, config.Get.StatusPort)
	defer rpc.StatusClose()

	_, err := net.Bootstrap(config.Get.ListenAddress, config.Get.ListenPort, config.Get.OnlineMode)
	if err != nil {
		panic(err) // TODO: tidy error reporting?
	}
}
