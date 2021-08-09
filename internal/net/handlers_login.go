package net

import (
	"errors"
	"github.com/JungleMC/java-edition/internal/net/packets"
	. "reflect"
)

func (c *JavaClient) loginHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(packets.ServerboundLoginStartPacket{}):
		return c.handleLoginStartPacket(pkt.(packets.ServerboundLoginStartPacket))
	}

	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleLoginStartPacket(pkt packets.ServerboundLoginStartPacket) error {
	return nil
}
