package main

import (
	"github.com/JungleMC/java-edition/pkg/startup"
	"github.com/JungleMC/sdk/pkg/redis"
)

func main() {
	startup.Start(redis.NewClient())
}
