package net

import (
	"errors"
	"github.com/JungleMC/java-edition/internal/config"
	"github.com/JungleMC/java-edition/internal/net/auth"
	"github.com/JungleMC/java-edition/internal/net/packets"
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
	/*
		msg := &messages.JavaLoginBegin{
			Username:            pkt.Username,
			GameProtocolVersion: c.gameProtocolVersion,
		}

		cmd := c.server.rdb.Publish(context.Background(), "login.java.begin", msg)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	*/

	c.authProfile.Name = pkt.Username

	if config.Get.OnlineMode {
		return c.send(&packets.ClientboundLoginEncryptionRequest{
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
		err := c.enableCompression()
		if err != nil {
			return err
		}
	}
	return nil
}
