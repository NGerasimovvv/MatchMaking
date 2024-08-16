package matchmaking

import (
	"encoding/json"
	"fmt"
	"github.com/NGerasimovvv/MatchMaking/internal/models"
	"github.com/NGerasimovvv/MatchMaking/internal/storage"
	"math"
	"net/http"
	"sort"
	"time"
)

type Matchmaker struct {
	storage   storage.Storage
	groupSize int
}

var groupNumber = 1

func NewMatchmaker(storage storage.Storage, groupSize int) *Matchmaker {
	return &Matchmaker{
		storage:   storage,
		groupSize: groupSize,
	}
}

func (m *Matchmaker) HandleAddPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player
	err := json.NewDecoder(r.Body).Decode(&player)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	player.JoinTime = time.Now().UTC()
	m.storage.AddPlayer(player)
	w.WriteHeader(http.StatusOK)

}

func (m *Matchmaker) FormGroups() {
	players := m.storage.GetAllPlayers()
	for len(players) >= m.groupSize {
		group := m.FindBestGroup(players, m.groupSize, 25, 25)
		if len(group.Players) == 0 {
			break
		}

		m.PrintGroupInfo(group)
		m.storage.RemovePlayers(group.Players)
		groupNumber++
		players = m.storage.GetAllPlayers()
	}
}

func (m *Matchmaker) FindBestGroup(players []models.Player, groupSize int, maxSkillDiff, maxLatencyDiff float64) models.Group {

	comb := combine(players, groupSize)
	bestGroup := models.Group{}
	bestTotalDifference := math.MaxFloat64
	for _, group := range comb {
		skillDiff, latencyDiff := CalculateGroupMaxDifferences(group)
		if skillDiff <= maxSkillDiff && latencyDiff <= maxLatencyDiff {
			if skillDiff+latencyDiff < bestTotalDifference {
				bestTotalDifference = skillDiff + latencyDiff
				bestGroup = models.Group{Players: group}
			}
		}
	}

	for _, group := range comb {
		if FilterGroup(group, 11.0, 11.0) {
			sort.Slice(group, func(i, j int) bool {
				return time.Since(group[i].JoinTime).Seconds() > time.Since(group[j].JoinTime).Seconds()
			})

			if len(group) > 0 {
				if len(bestGroup.Players) == 0 || len(group) > len(bestGroup.Players) {
					bestGroup = models.Group{Players: group}
				}
			}
		}
	}

	return bestGroup
}

func CalculateMaxDifference(values []float64) float64 {
	if len(values) < 2 {
		return 0.0
	}

	sort.Float64s(values)
	return values[len(values)-1] - values[0]
}

func CalculateGroupMaxDifferences(group []models.Player) (float64, float64) {
	var skills, latencies []float64

	for _, player := range group {
		skills = append(skills, player.Skill)
		latencies = append(latencies, player.Latency)
	}

	maxSkillDiff := CalculateMaxDifference(skills)
	maxLatencyDiff := CalculateMaxDifference(latencies)

	return maxSkillDiff, maxLatencyDiff
}

func FilterGroup(group []models.Player, maxSkillDiff, maxLatencyDiff float64) bool {
	var skills, latencies []float64

	for _, player := range group {
		skills = append(skills, player.Skill)
		latencies = append(latencies, player.Latency)
	}

	maxSkill := CalculateMaxDifference(skills)
	maxLatency := CalculateMaxDifference(latencies)

	return maxSkill <= maxSkillDiff && maxLatency <= maxLatencyDiff
}

func combine(players []models.Player, groupSize int) [][]models.Player {
	if groupSize == 0 {
		return [][]models.Player{{}}
	}
	if len(players) == 0 {
		return nil
	}

	var result [][]models.Player
	for i := 0; i <= len(players)-groupSize; i++ {
		for _, c := range combine(players[i+1:], groupSize-1) {
			result = append(result, append([]models.Player{players[i]}, c...))
		}
	}
	return result
}

func CalculateGroupStats(group []models.Player) (minSkill, maxSkill, avgSkill, minLatency, maxLatency, avgLatency, minWaitTime, maxWaitTime, avgWaitTime float64) {
	if len(group) == 0 {
		return
	}

	minSkill, maxSkill = math.MaxFloat64, -math.MaxFloat64
	minLatency, maxLatency = math.MaxFloat64, -math.MaxFloat64
	minWaitTime, maxWaitTime = math.MaxFloat64, -math.MaxFloat64

	var totalSkill, totalLatency, totalWaitTime float64

	for _, player := range group {
		if player.Skill < minSkill {
			minSkill = player.Skill
		}
		if player.Skill > maxSkill {
			maxSkill = player.Skill
		}
		totalSkill += player.Skill

		if player.Latency < minLatency {
			minLatency = player.Latency
		}
		if player.Latency > maxLatency {
			maxLatency = player.Latency
		}
		totalLatency += player.Latency

		waitTime := time.Since(player.JoinTime.UTC()).Seconds()
		if waitTime < minWaitTime {
			minWaitTime = waitTime
		}
		if waitTime > maxWaitTime {
			maxWaitTime = waitTime
		}
		totalWaitTime += waitTime
	}

	avgSkill = totalSkill / float64(len(group))
	avgLatency = totalLatency / float64(len(group))
	avgWaitTime = totalWaitTime / float64(len(group))

	return
}

func (m *Matchmaker) PrintGroupInfo(group models.Group) {
	minSkill, maxSkill, avgSkill, minLatency, maxLatency, avgLatency, minTime, maxTime, avgTime := CalculateGroupStats(group.Players)

	fmt.Printf("Group #%d:\n", groupNumber)
	fmt.Printf("  Min Skill: %.2f, Max Skill: %.2f, Avg Skill: %.2f\n", minSkill, maxSkill, avgSkill)
	fmt.Printf("  Min Latency: %.2f, Max Latency: %.2f, Avg Latency: %.2f\n", minLatency, maxLatency, avgLatency)
	fmt.Printf("  Min Wait Time: %.2fs, Max Wait Time: %.2fs, Avg Wait Time: %.2fs\n", minTime, maxTime, avgTime)
	fmt.Printf("  Players: ")
	for _, player := range group.Players {
		fmt.Printf("%s ", player.Name)
	}
	fmt.Println()
}
