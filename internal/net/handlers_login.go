package net

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net/auth"
	"github.com/JungleMC/java-edition/internal/net/packets"
	"github.com/JungleMC/sdk/pkg/events"
	"github.com/JungleMC/sdk/pkg/messages"
	"google.golang.org/protobuf/proto"
	. "reflect"
)

func (c *JavaClient) loginHandlers(pkt Packet) error {
	t := ValueOf(pkt).Type()
	switch t {
	case TypeOf(packets.ServerboundLoginStartPacket{}):
		return c.handleLoginStartPacket(pkt.(packets.ServerboundLoginStartPacket))
	case TypeOf(packets.ServerboundLoginEncryptionResponsePacket{}):
		return c.handleLoginEncryptionResponse(pkt.(packets.ServerboundLoginEncryptionResponsePacket))
	}

	return errors.New("not implemented: " + t.Name())
}

func (c *JavaClient) handleLoginStartPacket(pkt packets.ServerboundLoginStartPacket) error {
	c.authProfile.Name = pkt.Username

	if config.Get.OnlineMode {
		return c.Send(&packets.ClientboundLoginEncryptionRequest{
			ServerId:    "",
			PublicKey:   c.server.publicKeyBytes,
			VerifyToken: c.verifyToken,
		})
	} else {
		if config.Get.CompressionThreshold > 0 {
			err := c.enableCompression()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *JavaClient) handleLoginEncryptionResponse(pkt packets.ServerboundLoginEncryptionResponsePacket) error {
	sharedSecret, err := auth.DecryptLoginResponse(c.server.privateKey, c.server.publicKeyBytes, c.verifyToken, pkt.VerifyToken, pkt.SharedSecret, c.authProfile)
	if err != nil {
		return err
	}

	err = c.enableEncryption(sharedSecret)
	if err != nil {
		return err
	}

	if config.Get.CompressionThreshold > 0 {
		err = c.enableCompression()
		if err != nil {
			return err
		}
	}

	networkIdBytes, _ := c.networkId.MarshalBinary()
	profileIdBytes, _ := c.authProfile.ID.MarshalBinary()

	msg := &events.PlayerLoginEvent{
		ClientType: messages.ClientType_JAVA_EDITION,
		NetworkId:  networkIdBytes,
		ProfileId:  profileIdBytes,
		Username:   c.authProfile.Name,
	}

	err = c.setTextures(msg)
	if err != nil {
		return err
	}

	msgBytes, _ := proto.Marshal(msg)

	cmd := c.server.RDB.Publish(context.Background(), "event.login", msgBytes)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (c *JavaClient) setTextures(msg *events.PlayerLoginEvent) error {
	textures, err := c.decodeTextures()
	if err != nil {
		return err
	}

	if textures.Textures.Skin != nil {
		msg.SkinUrl = textures.Textures.Skin.Url
	}

	if textures.Textures.Cape != nil {
		msg.SkinUrl = textures.Textures.Cape.Url
	}
	return nil
}

func (c *JavaClient) decodeTextures() (*auth.TextureProperties, error) {
	textureJson, err := base64.StdEncoding.DecodeString(c.authProfile.Properties[0].Value)
	if err != nil {
		return nil, err
	}

	textures := &auth.TextureProperties{}
	err = json.Unmarshal(textureJson, textures)
	if err != nil {
		return nil, err
	}
	return textures, nil
}
