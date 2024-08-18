package dto

import (
	"matchMaker/internal/storage/postgres/repository/models"
)

type Group struct {
	Users       []models.User
	MinSkill    float64
	MaxSkill    float64
	AvgSkill    float64
	MinLatency  float64
	MaxLatency  float64
	AvgLatency  float64
	MinWaitTime float64
	MaxWaitTime float64
	AvgWaitTime float64
}
