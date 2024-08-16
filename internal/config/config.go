package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	GroupSize        int
	ServerAddress    string
	UseMemoryStorage bool
	DatabaseURL      string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port := os.Getenv("POSTGRES_PORT")
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	useMemoryStorage, _ := strconv.ParseBool(os.Getenv("UseMemoryStorage"))
	groupSize, _ := strconv.Atoi(os.Getenv("GroupSize"))
	serverAddress := os.Getenv("ServerAddress")

	databaseURL := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	return &Config{
		GroupSize:        groupSize,
		ServerAddress:    serverAddress,
		UseMemoryStorage: useMemoryStorage,
		DatabaseURL:      databaseURL,
	}
}
