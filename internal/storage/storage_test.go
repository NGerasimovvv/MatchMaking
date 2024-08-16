package storage

import (
	"github.com/NGerasimovvv/MatchMaking/internal/models"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func TestMemoryStorage(t *testing.T) {
	storage := NewMemoryStorage()

	player := models.Player{Name: "TestPlayer", Skill: 10.0, Latency: 50.0, JoinTime: time.Now().UTC()}
	storage.AddPlayer(player)

	players := storage.GetAllPlayers()
	if len(players) != 1 {
		t.Fatalf("Expected 1 player, got %d", len(players))
	}

	if players[0] != player {
		t.Errorf("Expected player %+v, got %+v", player, players[0])
	}

	storage.RemovePlayers(players)
	players = storage.GetAllPlayers()
	if len(players) != 0 {
		t.Fatalf("Expected 0 players, got %d", len(players))
	}
}

func TestDBStorage(t *testing.T) {

	databaseURL := "postgres://postgres:12345678@localhost:5432/matchmaking?sslmode=disable"
	db, err := NewDatabaseConnection(databaseURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	storage := NewDBStorage(db)

	player := models.Player{Name: "TestPlayer", Skill: 20.0, Latency: 30.0, JoinTime: time.Now().UTC()}
	storage.AddPlayer(player)

	players := storage.GetAllPlayers()
	if len(players) != 1 {
		t.Fatalf("Expected 1 player, got %d", len(players))
	}

	if players[0].Name != player.Name || players[0].Skill != player.Skill || players[0].Latency != player.Latency {
		t.Errorf("Expected player %+v, got %+v", player, players[0])
	}

	if players[0].JoinTime.Sub(player.JoinTime).Abs() > time.Second {
		t.Errorf("JoinTime mismatch: expected %v, got %v", player.JoinTime, players[0].JoinTime)
	}

	storage.RemovePlayers(players)
	players = storage.GetAllPlayers()
	if len(players) != 0 {
		t.Fatalf("Expected 0 players, got %d", len(players))
	}
}
