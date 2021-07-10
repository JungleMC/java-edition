package packets

import (
	"github.com/google/uuid"
)

type ServerListResponse struct {
	Description string    `json:"description,omitempty"` // chat message ptr
	Players     ServerListPlayers `json:"players"`
	Version     GameVersion       `json:"version"`
	Favicon     string            `json:"favicon,omitempty"`
}

type GameVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerListPlayers struct {
	Max    int                `json:"max"`
	Online int                `json:"online"`
	Sample []ServerListPlayer `json:"sample,omitempty"`
}

type ServerListPlayer struct {
	Name string    `json:"name"`
	Id   uuid.UUID `json:"id"`
}
