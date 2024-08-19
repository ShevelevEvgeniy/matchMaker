package models

import "time"

type User struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Skill           float64   `json:"skill"`
	Latency         float64   `json:"latency"`
	SearchMatch     bool      `json:"search_match"`
	SearchStartTime time.Time `json:"search_start_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
