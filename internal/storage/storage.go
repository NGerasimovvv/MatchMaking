package storage

import (
	"github.com/NGerasimovvv/MatchMaking/internal/models"
)

type Storage interface {
	AddPlayer(player models.Player)
	GetAllPlayers() []models.Player
	RemovePlayers(players []models.Player)
}
