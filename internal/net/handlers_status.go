package net

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/junglemc/Service-JavaEditionHost/internal/net/packets"
	player_msg "github.com/junglemc/Service-PlayerProvider/pkg/msg"
	player_rpc "github.com/junglemc/Service-PlayerProvider/pkg/rpc"
	status_msg "github.com/junglemc/Service-StatusProvider/pkg/msg"
	status_rpc "github.com/junglemc/Service-StatusProvider/pkg/rpc"
	. "reflect"
	"sync"
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
	statusResponse, statusResponseErr, playerListResponse, playerListResponseErr := getStatusInfo()
	if statusResponseErr != nil {
		return statusResponseErr
	}

	if playerListResponseErr != nil {
		return playerListResponseErr
	}

	players := make([]packets.ServerListPlayer, len(playerListResponse.Sample))
	for i := 0; i < len(players); i++ {
		id, _ := uuid.FromBytes(playerListResponse.Sample[i].Id)
		players[i] = packets.ServerListPlayer{
			Name: playerListResponse.Sample[i].Name,
			Id:   id,
		}
	}

	status := packets.ServerListResponse{
		Description: statusResponse.ServerDescription,
		Players: packets.ServerListPlayers{
			Max:    playerListResponse.Max,
			Online: playerListResponse.Online,
			Sample: []packets.ServerListPlayer{},
		},
		Version: packets.GameVersion{
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

// TODO: Better handling of this data, it's a mess, but it's parallel
func getStatusInfo() (statusResponse *status_msg.StatusResponse, statusResponseErr error, playerListResponse *player_msg.JavaEdition_PlayerListResponse, playerListResponseErr error) {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		statusResponse, statusResponseErr = status_rpc.StatusService.StatusRequest(ctx, &status_msg.StatusRequest{})
		if statusResponseErr != nil {
			status_rpc.StatusServiceConnection.ResetConnectBackoff()
		}
	}()

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		playerListResponse, playerListResponseErr = player_rpc.PlayerService.JavaEdition_PlayerListRequest(ctx, &player_msg.JavaEdition_PlayerListRequest{})
		if playerListResponseErr != nil {
			player_rpc.PlayerServiceConnection.ResetConnectBackoff()
		}
	}()

	wg.Wait()
	return
}