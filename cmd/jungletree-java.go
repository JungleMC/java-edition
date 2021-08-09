package main

import (
	"github.com/JungleMC/java-edition/pkg/service"
	"github.com/JungleMC/sdk/pkg/redis"
)

func main() {
	_ = service.Start(redis.NewClient())
}
