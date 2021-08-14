package main

import (
	"bytes"
	"github.com/JungleMC/java-edition/internal/net"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
)

func testData(filename string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata/", filename))
	if err != nil {
		panic(err)
	}
	return data
}

func main() {
	data := testData("handshake/handshake_packet.bin")
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()

	for i:=0; i<1000000; i++ {
		_, err := net.ReadPacket(bytes.NewBuffer(data), net.Handshake, false)
		if err != nil {
			panic(err)
		}
	}
}
