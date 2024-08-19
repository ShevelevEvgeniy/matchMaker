package user_converter

import (
	"time"

	clust "github.com/muesli/clusters"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

func ServiceToRepoModel(user dto.User) models.User {
	return models.User{
		Name:            user.Name,
		Skill:           user.Skill,
		Latency:         user.Latency,
		SearchMatch:     true,
		SearchStartTime: time.Now(),
	}
}

func UsersToIds(users []models.User) []int64 {
	IDs := make([]int64, len(users))

	for i, user := range users {
		IDs[i] = user.ID
	}

	return IDs
}

func UsersToMatrix(users []models.User) ([]clust.Coordinates, map[int64]int) {
	var dataset []clust.Coordinates
	userIndexMap := make(map[int64]int)

	for i, user := range users {
		obs := clust.Coordinates{
			user.Skill,
			user.Latency,
		}
		dataset = append(dataset, obs)

		userIndexMap[user.ID] = i
	}

	return dataset, userIndexMap
}
