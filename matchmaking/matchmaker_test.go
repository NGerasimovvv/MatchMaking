package matchmaking

import (
	"bytes"
	"encoding/json"
	"github.com/NGerasimovvv/MatchMaking/internal/models"
	"github.com/NGerasimovvv/MatchMaking/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleAddPlayer(t *testing.T) {
	// Создаем временное хранилище для теста
	memStorage := storage.NewMemoryStorage()

	// Создаем матчмейкер с временным хранилищем
	matchmaker := NewMatchmaker(memStorage, 5)

	// Создаем данные игрока
	player := models.Player{
		Name:    "TestPlayer",
		Skill:   40,
		Latency: 10,
	}

	// Преобразуем игрока в JSON
	playerJSON, err := json.Marshal(player)
	if err != nil {
		t.Fatalf("Failed to marshal player: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(playerJSON))
	rr := httptest.NewRecorder()

	matchmaker.HandleAddPlayer(rr, req)

	// Проверяем статус ответа
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	// Проверяем, что игрок был добавлен
	players := memStorage.GetAllPlayers()
	if len(players) != 1 {
		t.Errorf("Expected 1 player, got %d", len(players))
	}

	// Проверяем данные игрока
	storedPlayer := players[0]
	if storedPlayer.Name != player.Name {
		t.Errorf("Expected Name %s, got %s", player.Name, storedPlayer.Name)
	}
	if storedPlayer.Skill != player.Skill {
		t.Errorf("Expected Skill %f, got %f", player.Skill, storedPlayer.Skill)
	}
	if storedPlayer.Latency != player.Latency {
		t.Errorf("Expected Latency %f, got %f", player.Latency, storedPlayer.Latency)
	}

	// Проверяем, что разница во времени не слишком велика
	now := time.Now().UTC()
	if storedPlayer.JoinTime.Before(now.Add(-time.Second)) || storedPlayer.JoinTime.After(now.Add(time.Second)) {
		t.Errorf("JoinTime %v is out of expected range [%v, %v]", storedPlayer.JoinTime, now.Add(-time.Second), now.Add(time.Second))
	}
}

func TestFindBestGroup(t *testing.T) {
	storage := storage.NewMemoryStorage()

	players := []models.Player{
		{Name: "Player1", Skill: 10, Latency: 20, JoinTime: time.Now().UTC()},
		{Name: "Player2", Skill: 15, Latency: 25, JoinTime: time.Now().UTC()},
		{Name: "Player3", Skill: 20, Latency: 30, JoinTime: time.Now().UTC()},
		{Name: "Player4", Skill: 25, Latency: 35, JoinTime: time.Now().UTC()},
		{Name: "Player5", Skill: 30, Latency: 40, JoinTime: time.Now().UTC()},
		{Name: "Player6", Skill: 35, Latency: 45, JoinTime: time.Now().UTC()},
	}

	for _, player := range players {
		storage.AddPlayer(player)
	}

	matchmaker := NewMatchmaker(storage, 5)
	group := matchmaker.FindBestGroup(players, 5, 20, 20)

	if len(group.Players) != 5 {
		t.Fatalf("Expected group size 5, got %d", len(group.Players))
	}
}
