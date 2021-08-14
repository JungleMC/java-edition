package net

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/pprof"
	"testing"
)

func testData(filename string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata/", filename))
	if err != nil {
		panic(err)
	}
	return data
}

func BenchmarkReadPacket(b *testing.B) {
	data := testData("handshake/handshake_packet.bin")
	pprof.StartCPUProfile(os.Stdout)
	for i := 0; i < b.N; i++ {
		_, err := ReadPacket(bytes.NewBuffer(data), Handshake, false)
		if err != nil {
			panic("uh")
		}
	}
	pprof.StopCPUProfile()
}
