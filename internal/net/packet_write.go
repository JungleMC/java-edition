package net

import (
	"bytes"
	"compress/zlib"
	"github.com/JungleMC/java-edition/internal/net/codec"
	"reflect"
)

func WritePacket(buf *bytes.Buffer, v reflect.Value, proto Protocol, compressed bool, compressionThreshold int) {
	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	id := packetIDFromTypeClientbound(v.Type(), proto)

	packet := append(codec.WriteVarInt32(id), codec.Marshal(v.Interface())...)

	if compressed {
		if len(packet) >= compressionThreshold {
			packet = compress(packet)
		} else {
			packet = append(codec.WriteVarInt32(0), packet...)
		}
	}

	buf.Write(codec.WriteVarInt32(int32(len(packet))))
	buf.Write(packet)
}

func compress(data []byte) []byte {
	buf := &bytes.Buffer{}
	writer := zlib.NewWriter(buf)
	_, _ = writer.Write(data)
	_ = writer.Flush()
	return append(codec.WriteVarInt32(int32(len(data))), buf.Bytes()...)
}
