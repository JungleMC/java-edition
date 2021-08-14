package main

import (
	"bytes"
	"github.com/JungleMC/java-edition/internal/net"
	"github.com/JungleMC/java-edition/internal/net/packets"
	"os"
	"reflect"
)

func main() {
	statusRequest := packets.ServerboundHandshakePacket{
		ProtocolVersion: 310,
		ServerHost:      "localhost",
		ServerPort:      25565,
		NextState:       int32(net.Status),
	}

	buf := &bytes.Buffer{}
	net.WritePacket(buf, reflect.ValueOf(statusRequest), net.Handshake, false, 0)
	err := os.WriteFile("testdata/handshake_packet.bin", buf.Bytes(), 0o755)
	if err != nil {
		panic(err)
	}
}
