package service

import (
	"fmt"
	"github.com/JungleMC/protocol"
	"github.com/JungleMC/sdk/pkg/events"
	"github.com/JungleMC/sdk/pkg/messages"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (s *JavaService) listenLogin(loginChannel <-chan *redis.Message) error {
	for msg := range loginChannel {
		response := &events.PlayerLoginResponse{}
		err := messages.UnmarshalMessage(msg, proto.Message(response))
		if err != nil {
			return err
		}

		if response.GetUsername() == "" {
			return fmt.Errorf("username not returned to the java host")
		}

		client, err := s.NetworkServer.GetClient(response.GetNetworkId())
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}

		profileId, err := uuid.FromBytes(response.GetProfileId())
		if err != nil {
			return fmt.Errorf("failed to read profile ID: %w", err)
		}

		err = client.Send(&protocol.LoginSuccess{
			UUID:     profileId,
			Username: response.GetUsername(),
		})
		if err != nil {

		}
	}
	return nil
}
