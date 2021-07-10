package main

import (
	"github.com/junglemc/Service-JavaEditionHost/internal/net"
	"os"
	"strconv"
)

func main() {
	javaAddr, javaPort := javaListenAddress()

	onlineModeEnv := os.Getenv("JAVA_ONLINE_MODE")
	onlineMode, err := strconv.ParseBool(onlineModeEnv)
	if err != nil {
		panic("JAVA_ONLINE_MODE not bool: " + onlineModeEnv)
	}

	_, err = net.Bootstrap(javaAddr, javaPort, onlineMode)
}

func javaListenAddress() (string, int) {
	javaAddrEnv := os.Getenv("JAVA_LISTEN_ADDRESS")
	javaPortEnv := os.Getenv("JAVA_LISTEN_PORT")

	if javaAddrEnv == "" {
		javaAddrEnv = "0.0.0.0"
	}

	if javaPortEnv == "" {
		javaPortEnv = "25565"
	}
	javaPort, err := strconv.Atoi(javaPortEnv)
	if err != nil {
		panic("JAVA_LISTEN_PORT invalid: " + javaPortEnv)
	}
	return javaAddrEnv, javaPort
}
