package storage

import (
	"github.com/NGerasimovvv/MatchMaking/internal/models"
	"sync"
)

type MemoryStorage struct {
	mu      sync.Mutex
	players []models.Player
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		players: make([]models.Player, 0),
	}
}

func (s *MemoryStorage) AddPlayer(player models.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.players = append(s.players, player)
}

func (s *MemoryStorage) GetAllPlayers() []models.Player {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.players
}

func (s *MemoryStorage) RemovePlayers(players []models.Player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	groupMap := make(map[string]bool)
	for _, player := range players {
		groupMap[player.Name] = true
	}

	var remainingPlayers []models.Player
	for _, player := range s.players {
		if !groupMap[player.Name] {
			remainingPlayers = append(remainingPlayers, player)
		}
	}
	s.players = remainingPlayers
}
