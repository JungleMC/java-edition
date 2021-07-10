package net

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/junglemc/Service-JavaEditionHost/internal/net/packets"
	"github.com/junglemc/Service-JavaEditionHost/internal/rpc"
	"github.com/junglemc/Service-StatusProvider/pkg/msg"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO: Query player tracking service in parallel with the status service
	statusResponse, err := rpc.Status.StatusRequest(ctx, &msg.StatusRequest{})
	if err != nil {
		rpc.StatusConnection.ResetConnectBackoff()
		return err
	}

	status := packets.ServerListResponse{
		Description: statusResponse.ServerDescription,
		Players:     packets.ServerListPlayers{
			Max:    20, // TODO: Query player tracking service
			Online: 0, // TODO: Query player tracking service
			Sample: []packets.ServerListPlayer{}, // TODO: Populate online player list sample from player tracking service
		},
		Version:     packets.GameVersion{
			Name:     ProtocolVersionName,
			Protocol: ProtocolVersionCode,
		},
		Favicon: statusResponse.GetFavicon(),
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