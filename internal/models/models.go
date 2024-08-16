package models

import "time"

type Player struct {
	Name     string
	Skill    float64
	Latency  float64
	JoinTime time.Time
}

type Group struct {
	Players []Player
}
