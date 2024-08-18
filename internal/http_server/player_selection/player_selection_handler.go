package player_selection_handler

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/muesli/clusters"
	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/storage/postgres/repository/models"
)

type Service interface {
	GetUsersInSearch(ctx context.Context, groupSize int) ([]models.User, error)
	SaveRemainingUsers(ctx context.Context, users []models.User) error
	GetAndRemoveRemainingUsers(ctx context.Context) ([]models.User, bool, error)
}

type LogGroups interface {
	PrintGroupInfo(ctx context.Context, groups []dto.Group)
}

type PlayerSelectionInterface interface {
	Run(ctx context.Context)
}

type PlayerSelection struct {
	cfg       config.MatchSettings
	log       *zap.Logger
	service   Service
	logGroups LogGroups
}

func NewPlayerSelection(
	cfg config.MatchSettings,
	log *zap.Logger,
	service Service,
	logGroups LogGroups,
) *PlayerSelection {

	return &PlayerSelection{
		cfg:       cfg,
		log:       log,
		service:   service,
		logGroups: logGroups,
	}
}

func (p *PlayerSelection) delay() {
	time.Sleep(p.cfg.Delay * time.Second)
}

func (p *PlayerSelection) Run(ctx context.Context) {
	p.log.Info("player selection started")

	usersChan := p.selectUsers(ctx)
	go p.generateGroups(ctx, usersChan)
}

func (p *PlayerSelection) selectUsers(ctx context.Context) <-chan []models.User {
	usersChan := make(chan []models.User, 1)

	go func() {
		defer func() {
			close(usersChan)
		}()

		for {
			select {
			case <-ctx.Done():
				p.log.Info("context done, skipping user selection")
				return
			default:
				users, err := p.service.GetUsersInSearch(ctx, p.cfg.BatchSize)
				if err != nil {
					p.log.Error("error occurred on getting users in search:", zap.String("error", err.Error()))
					p.delay()
				}

				usersChan <- users
			}
		}
	}()

	return usersChan
}

func (p *PlayerSelection) generateGroups(ctx context.Context, usersChan <-chan []models.User) {
	for {
		select {
		case <-ctx.Done():
			p.log.Info("context done, skipping group generation")
			return
		case users, ok := <-usersChan:
			if !ok {
				p.log.Info("users channel closed, skipping group generation")
				return
			}

			cachedUsers, hasCachedUsers, cacheErr := p.service.GetAndRemoveRemainingUsers(ctx)
			if cacheErr != nil {
				p.log.Error("error occurred on getting cached users:", zap.String("error", cacheErr.Error()))
			}

			if len(users) == 0 {
				p.log.Info("no users in search")
				p.delay()
				continue
			}

			if hasCachedUsers {
				users = append(cachedUsers, users...)
			}

			groups := p.createGroupsUsingNearestNeighbors(ctx, users)
			if len(groups) > 0 {
				p.logGroups.PrintGroupInfo(ctx, groups)
			}
		}
	}
}

func (p *PlayerSelection) createGroupsUsingNearestNeighbors(ctx context.Context, users []models.User) []dto.Group {
	userMatrix, userIndexMap := converter.UsersToMatrix(users)

	var groups []dto.Group
	var remainingUsers []models.User
	groupedUsers := make(map[int]bool)

	for i := 0; i < len(users); i++ {
		if groupedUsers[users[i].ID] {
			continue
		}

		currentGroup := dto.Group{Users: []models.User{}}
		closestUsers := p.findClosestUsers(userMatrix, userIndexMap[users[i].ID])

		for _, closest := range closestUsers {
			if !groupedUsers[closest.Index] && p.isWithinSkillRange(users[closest.Index].Skill, users[i].Skill) {
				groupedUsers[closest.Index] = true
				currentGroup.Users = append(currentGroup.Users, users[closest.Index])
			}
		}

		if len(currentGroup.Users) == p.cfg.GroupSize {
			groups = append(groups, currentGroup)
		} else {
			for _, user := range currentGroup.Users {
				remainingUsers = append(remainingUsers, user)
			}
		}
	}

	if len(remainingUsers) > 0 {
		if err := p.service.SaveRemainingUsers(ctx, remainingUsers); err != nil {
			p.log.Error("error occurred on saving remaining users:", zap.String("error", err.Error()))
		}
	}

	return groups
}

func (p *PlayerSelection) findClosestUsers(userMatrix []clusters.Coordinates, index int) []dto.UserDistance {
	var distances []dto.UserDistance
	for j, userCoordinates := range userMatrix {
		if j == index {
			continue
		}
		distance := euclideanDistance(userMatrix[index], userCoordinates)
		distances = append(distances, dto.UserDistance{Index: j, Distance: distance})
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Distance < distances[j].Distance
	})

	var closestUsers []dto.UserDistance
	for _, dist := range distances {
		if len(closestUsers) >= p.cfg.GroupSize {
			break
		}
		closestUsers = append(closestUsers, dist)
	}

	return closestUsers
}

func euclideanDistance(a, b clusters.Coordinates) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

func (p *PlayerSelection) isWithinSkillRange(skill1, skill2 float64) bool {
	return math.Abs(skill1-skill2) <= p.cfg.RangeSkill
}
