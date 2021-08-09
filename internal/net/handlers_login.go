package net

import (
	"context"
	"errors"
	"github.com/JungleMC/java-edition/internal/net/packets"
	"github.com/JungleMC/sdk/pkg/redis/messages"
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
	msg := &messages.JavaLoginBegin{
		Username: pkt.Username,
		GameProtocolVersion: c.gameProtocolVersion,
	}

	cmd := c.server.RDB.Publish(context.Background(), "login.java.begin", msg)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}
