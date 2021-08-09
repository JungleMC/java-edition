package main

import (
	"github.com/JungleMC/java-edition/pkg/service"
	"github.com/JungleMC/sdk/pkg/redis"
)

func main() {
	service.Start(redis.NewClient())
}
