package net

import (
	"errors"
	"github.com/junglemc/Service-JavaEditionHost/internal/net/packets"
	. "reflect"
)

func (c *JavaClient) handshakeHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(packets.ServerboundHandshakePacket{}):
		return c.handleHandshakePacket(pkt.(packets.ServerboundHandshakePacket))
	case TypeOf(packets.ServerboundHandshakeLegacyPingPacket{}):
		return c.handleHandshakeLegacyPing(pkt.(packets.ServerboundHandshakeLegacyPingPacket))
	}
	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleHandshakePacket(pkt packets.ServerboundHandshakePacket) error {
	c.protocol = Protocol(pkt.NextState)
	c.gameProtocolVersion = pkt.ProtocolVersion
	return nil
}

func (c *JavaClient) handleHandshakeLegacyPing(pkt packets.ServerboundHandshakeLegacyPingPacket) error {
	c.protocol = Status
	c.gameProtocolVersion = int32(pkt.ProtocolVersion)
	return nil
}
