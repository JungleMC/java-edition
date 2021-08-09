package net

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/JungleMC/java-edition/internal/net/packets"
	. "reflect"
	"time"
)

func (c *JavaClient) statusHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(packets.ServerboundStatusRequestPacket{}):
		return c.handleStatusRequest()
	case TypeOf(packets.ServerboundStatusPingPacket{}):
		return c.handleStatusPing()
	}
	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleStatusRequest() error {
	description := c.server.rdb.Get(context.Background(), "config:description").Val()
	favicon := c.server.rdb.Get(context.Background(), "config:favicon").Val()
	maxPlayers, _ := c.server.rdb.Get(context.Background(), "config:max_players").Int()

	status := packets.ServerListResponse{
		Description: description,
		Players:     packets.ServerListPlayers{
			Max:    maxPlayers,
			Online: 0,
			Sample: make([]packets.ServerListPlayer, 0),
		},
		Version: packets.GameVersion{
			Name:     ProtocolVersionName,
			Protocol: ProtocolVersionCode,
		},
		Favicon: favicon,
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return c.send(&packets.ClientboundStatusResponsePacket{Response: string(data)})
}

func (c *JavaClient) handleStatusPing() error {
	return c.send(&packets.ClientboundStatusPongPacket{Time: time.Now().Unix()})
}
