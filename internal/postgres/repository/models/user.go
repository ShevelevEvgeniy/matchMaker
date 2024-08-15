package models

import "time"

type User struct {
	Name            string    `json:"Name"`
	Skill           float32   `json:"Skill"`
	Latency         float32   `json:"Latency"`
	SearchingMatch  bool      `json:"SearchingMatch"`
	SearchStartTime time.Time `json:"SearchStartTime"`
	CreatedAt       time.Time `json:"CreatedAt"`
	UpdatedAt       time.Time `json:"UpdatedAt"`
}
