package main

import (
	"github.com/NGerasimovvv/MatchMaking/internal/config"
	"github.com/NGerasimovvv/MatchMaking/internal/storage"
	"github.com/NGerasimovvv/MatchMaking/matchmaking"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

func main() {

	cfg := config.LoadConfig()

	var currentStorage storage.Storage
	if cfg.UseMemoryStorage {
		currentStorage = storage.NewMemoryStorage()
	} else {
		db, err := storage.NewDatabaseConnection(cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
		currentStorage = storage.NewDBStorage(db)
	}

	matchmaker := matchmaking.NewMatchmaker(currentStorage, cfg.GroupSize)

	go func() {
		for {
			matchmaker.FormGroups()
			time.Sleep(1 * time.Second)
		}
	}()

	http.HandleFunc("/users", matchmaker.HandleAddPlayer)
	log.Println("Matchmaker service started")
	err := http.ListenAndServe(cfg.ServerAddress, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
