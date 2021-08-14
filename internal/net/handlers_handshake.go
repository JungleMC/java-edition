package net

import (
	"errors"
	. "github.com/JungleMC/protocol"
	. "reflect"
)

func (c *JavaClient) handshakeHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(&HandshakePacket{}):
		return c.handleHandshakePacket(pkt.(*HandshakePacket))
	case TypeOf(&LegacyPingPacket{}):
		return c.handleHandshakeLegacyPing(pkt.(*LegacyPingPacket))
	}
	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleHandshakePacket(pkt *HandshakePacket) error {
	c.state = ConnectionState(pkt.NextState)
	c.gameProtocolVersion = pkt.ProtocolVersion
	return nil
}

func (c *JavaClient) handleHandshakeLegacyPing(pkt *LegacyPingPacket) error {
	c.state = Status
	c.gameProtocolVersion = int32(pkt.ProtocolVersion)
	return nil
}
