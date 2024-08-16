package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Player struct {
	Name    string  `json:"name"`
	Skill   float64 `json:"skill"`
	Latency float64 `json:"latency"`
}

func main() {
	players := []Player{
		{Name: "player1", Skill: 20.0, Latency: 30.0},
		{Name: "player2", Skill: 25.0, Latency: 35.0},
		{Name: "player3", Skill: 30.0, Latency: 40.0},
		{Name: "player4", Skill: 35.0, Latency: 45.0},
		{Name: "player5", Skill: 40.0, Latency: 50.0},
		{Name: "player6", Skill: 45.0, Latency: 55.0},
		{Name: "player7", Skill: 50.0, Latency: 60.0},
		{Name: "player8", Skill: 55.0, Latency: 65.0},
		{Name: "player9", Skill: 60.0, Latency: 70.0},
		{Name: "player10", Skill: 65.0, Latency: 75.0},
		{Name: "player11", Skill: 65.0, Latency: 75.0},
		{Name: "player12", Skill: 30.0, Latency: 40.0},
	}

	for _, player := range players {
		go func(p Player) {
			jsonData, err := json.Marshal(p)
			if err != nil {
				fmt.Printf("Error marshalling player data: %v\n", err)
				return
			}
			resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Printf("Error sending POST request: %v\n", err)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("Player %s added with status code %d\n", p.Name, resp.StatusCode)
		}(player)

		time.Sleep(100 * time.Millisecond)
	}
}
