package packets

import (
	"github.com/google/uuid"
)

type ServerListResponse struct {
	Description string            `json:"description,omitempty"` // chat message ptr
	Players     ServerListPlayers `json:"players"`
	Version     GameVersion       `json:"version"`
	Favicon     string            `json:"favicon,omitempty"`
}

type GameVersion struct {
	Name     string `json:"name"`
	Protocol uint32 `json:"protocol"`
}

type ServerListPlayers struct {
	Max    uint32             `json:"max"`
	Online uint32             `json:"online"`
	Sample []ServerListPlayer `json:"sample,omitempty"`
}

type ServerListPlayer struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}
