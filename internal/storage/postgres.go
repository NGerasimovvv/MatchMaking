package storage

import (
	"database/sql"
	"github.com/NGerasimovvv/MatchMaking/internal/models"
	"log"
	"time"
)

type DBStorage struct {
	db *sql.DB
}

func NewDBStorage(db *sql.DB) *DBStorage {
	return &DBStorage{db: db}
}

func (s *DBStorage) AddPlayer(player models.Player) {
	_, err := s.db.Exec("INSERT INTO players (name, skill, latency, join_time) VALUES ($1, $2, $3, $4)",
		player.Name, player.Skill, player.Latency, player.JoinTime)
	if err != nil {
		log.Printf("Failed to add player: %v", err)
	}
}

func (s *DBStorage) GetAllPlayers() []models.Player {
	rows, err := s.db.Query("SELECT name, skill, latency, join_time FROM players")
	if err != nil {
		log.Printf("Failed to get players: %v", err)
		return nil
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		var joinTime time.Time
		err := rows.Scan(&player.Name, &player.Skill, &player.Latency, &joinTime)
		if err != nil {
			log.Printf("Failed to scan player: %v", err)
			continue
		}
		player.JoinTime = joinTime.UTC()
		players = append(players, player)
	}
	return players
}

func (s *DBStorage) RemovePlayers(players []models.Player) {
	for _, player := range players {
		_, err := s.db.Exec("DELETE FROM players WHERE name = $1", player.Name)
		if err != nil {
			log.Printf("Failed to remove player: %v", err)
		}
	}
}

func NewDatabaseConnection(databaseURL string) (*sql.DB, error) {
	return sql.Open("postgres", databaseURL)
}
