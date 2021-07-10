package main

import (
	"github.com/caarlos0/env"
	"github.com/junglemc/Service-JavaEditionHost/internal/config"
	"github.com/junglemc/Service-JavaEditionHost/internal/net"
	player_rpc "github.com/junglemc/Service-PlayerProvider/pkg/rpc"
	status_rpc "github.com/junglemc/Service-StatusProvider/pkg/rpc"
)

func main() {
	config.Get = &config.Config{}
	if err := env.Parse(config.Get); err != nil {
		panic(err)
	}

	status_rpc.StatusInit(config.Get.StatusHost, config.Get.StatusPort)
	defer status_rpc.StatusClose()

	player_rpc.PlayerInit(config.Get.PlayerHost, config.Get.PlayerPort)
	defer player_rpc.PlayerClose()

	_, err := net.Bootstrap(config.Get.ListenAddress, config.Get.ListenPort, config.Get.OnlineMode)
	if err != nil {
		panic(err) // TODO: tidy error reporting?
	}
}
