package converter

import (
	"time"

	"matchMaker/internal/dto"
	"matchMaker/internal/postgres/repository/models"
)

func ServiceToRepoModel(user dto.User) models.User {
	return models.User{
		Name:            user.Name,
		Skill:           user.Skill,
		Latency:         user.Latency,
		SearchingMatch:  true,
		SearchStartTime: time.Now(),
	}
}
