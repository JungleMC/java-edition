package net

import (
	"context"
	"encoding/json"
	"errors"
	. "github.com/JungleMC/protocol"
	. "reflect"
	"time"
)

func (c *JavaClient) statusHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(&StatusRequest{}):
		return c.handleStatusRequest()
	case TypeOf(&StatusPing{}):
		return c.handleStatusPing()
	}
	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleStatusRequest() error {
	description := c.server.RDB.Get(context.Background(), "config:description").Val()
	favicon := c.server.RDB.Get(context.Background(), "config:favicon").Val()
	maxPlayers, _ := c.server.RDB.Get(context.Background(), "config:max_players").Int()

	status := ServerListResponse{
		Description: description,
		Players:     ServerListPlayers{
			Max:    maxPlayers,
			Online: 0,
			Sample: make([]ServerListPlayer, 0),
		},
		Version: GameVersion{
			Name:     VersionDescription,
			Protocol: Version,
		},
		Favicon: favicon,
	}

	data, err := json.Marshal(status)
	if err != nil {
		return err
	}

	return c.Send(&StatusResponse{Response: string(data)})
}

func (c *JavaClient) handleStatusPing() error {
	return c.Send(&StatusPong{Time: time.Now().Unix()})
}
