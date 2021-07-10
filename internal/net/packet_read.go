package net

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"github.com/junglemc/Service-JavaEditionHost/internal/net/codec"
	"reflect"
)

type Packet interface{}

func readPacket(buf *bytes.Buffer, proto Protocol, compressed bool) (Packet, error) {
	payloadCheck, err := buf.ReadByte()
	if payloadCheck == 0xFE {
		readLegacyPing(buf)
		return nil, nil
	} else {
		_ = buf.UnreadByte()
	}

	var uncompressedLength int32
	var reader *bufio.Reader
	if compressed {
		compressedLength := codec.ReadVarInt32(buf)
		uncompressedLength = codec.ReadVarInt32(buf)

		if uncompressedLength > 0 {
			data := make([]byte, compressedLength)
			_, _ = buf.Read(data)
			zlibReader, err := zlib.NewReader(bytes.NewBuffer(data))
			if err != nil {
				return nil, err
			}
			reader = bufio.NewReader(zlibReader)
		} else {
			reader = bufio.NewReader(buf)
		}
	} else {
		uncompressedLength = codec.ReadVarInt32(buf)
		if err != nil {
			return nil, err
		}
		reader = bufio.NewReader(buf)
	}

	data := make([]byte, uncompressedLength)
	_, err = reader.Read(data)
	if err != nil {
		return nil, err
	}

	// Redefine the bytes reader here
	buf = bytes.NewBuffer(data)
	id := codec.ReadVarInt32(buf)

	pktType := packetTypeFromIDServerbound(id, proto)
	if pktType == nil {
		panic("nil type")
	}

	pkt := reflect.New(pktType).Elem()
	err = codec.Unmarshal(buf.Bytes(), pkt)
	if err != nil {
		return nil, err
	}
	return pkt.Interface().(Packet), err
}

func readLegacyPing(buf *bytes.Buffer) {
	_, _ = buf.ReadByte()
	_, _ = buf.ReadByte()
	_, _ = buf.ReadByte() // packet identifier for a plugin message

	mcPingHostLength := codec.ReadUint16(buf)
	mcPingHost := make([]byte, mcPingHostLength)
	_, _ = buf.Read(mcPingHost)

	codec.ReadInt16(buf) // Remaining
	_, _ = buf.ReadByte()

	hostnameLength := codec.ReadInt16(buf)
	hostname := make([]byte, hostnameLength)
	_, _ = buf.Read(hostname)
	codec.ReadInt16(buf)
}
