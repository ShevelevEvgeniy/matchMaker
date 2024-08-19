package dto

import (
	"time"

	"matchMaker/internal/storage/postgres/repository/models"
)

type Group struct {
	Users               []models.User
	GroupID             int
	MinSkill            float64
	MaxSkill            float64
	AvgSkill            float64
	MinLatency          float64
	MaxLatency          float64
	AvgLatency          float64
	MinWaitTime         float64
	MaxWaitTime         float64
	AvgWaitTime         float64
	MinTimeSpentInQueue time.Duration
	MaxTimeSpentInQueue time.Duration
	AvgTimeSpentInQueue time.Duration
}
