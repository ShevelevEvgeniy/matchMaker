package player_selection_handler

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/muesli/clusters"
	"go.uber.org/zap"
	"matchMaker/config"
	userConverter "matchMaker/internal/converter/user_converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/http_server/events"
	"matchMaker/internal/storage/postgres/repository/models"
)

type Service interface {
	GetUsersInSearch(ctx context.Context, groupSize int) ([]models.User, error)
	SaveRemainingUsers(ctx context.Context, users []models.User) error
	GetAndRemoveRemainingUsers(ctx context.Context) ([]models.User, bool, error)
}

type Events interface {
	Handle(ctx context.Context, message events.Message)
}

type GroupFormationHandler struct {
	cfg          config.MatchSettings
	log          *zap.Logger
	service      Service
	events       Events
	groupCounter int
	mu           sync.Mutex
}

func NewPlayerSelection(
	cfg config.MatchSettings,
	log *zap.Logger,
	service Service,
	events Events,
) *GroupFormationHandler {

	return &GroupFormationHandler{
		cfg:     cfg,
		log:     log,
		service: service,
		events:  events,
	}
}

func (g *GroupFormationHandler) delay() {
	time.Sleep(g.cfg.Delay * time.Second)
}

func (g *GroupFormationHandler) Run(ctx context.Context) {
	g.log.Info("player selection started")

	usersChan := g.selectUsers(ctx)

	for i := 0; i < g.cfg.CountWorkers; i++ {
		go g.generateGroups(ctx, usersChan)
	}
}

func (g *GroupFormationHandler) selectUsers(ctx context.Context) <-chan []models.User {
	usersChan := make(chan []models.User, g.cfg.CountWorkers)

	go func() {
		defer close(usersChan)

		for {
			select {
			case <-ctx.Done():
				g.log.Error("context done, skipping group generation", zap.String("error", ctx.Err().Error()))
				return
			default:
				users, err := g.service.GetUsersInSearch(ctx, g.cfg.BatchSize)
				if err != nil {
					g.log.Error("error occurred on getting users in search:", zap.String("error", err.Error()))
				}

				if len(users) == 0 {
					g.delay()
					continue
				}

				usersChan <- users
			}
		}
	}()

	return usersChan
}

func (g *GroupFormationHandler) generateGroups(ctx context.Context, usersChan <-chan []models.User) {
	for {
		select {
		case <-ctx.Done():
			g.log.Error("context done, skipping group generation", zap.String("error", ctx.Err().Error()))
			return
		case users, ok := <-usersChan:
			if !ok {
				g.log.Info("users channel closed, skipping group generation")
				return
			}

			cachedUsers, hasCachedUsers, cacheErr := g.service.GetAndRemoveRemainingUsers(ctx)
			if cacheErr != nil {
				g.log.Error("error occurred on getting cached users:", zap.String("error", cacheErr.Error()))
			}

			if hasCachedUsers {
				users = append(cachedUsers, users...)
			}

			g.createGroupsUsingNearestNeighbors(ctx, users)
		}
	}
}

func (g *GroupFormationHandler) createGroupsUsingNearestNeighbors(ctx context.Context, users []models.User) {
	userMatrix, userIndexMap := userConverter.UsersToMatrix(users)

	var remainingUsers []models.User
	groupedUsers := make(map[int]struct{})

	for i := 0; i < len(users); i++ {
		if _, exists := groupedUsers[users[i].ID]; exists {
			continue
		}

		closestUsers := g.findClosestUsers(userMatrix, userIndexMap[users[i].ID], groupedUsers)
		currentGroup := g.formGroup(users, closestUsers, groupedUsers)

		if len(currentGroup.Users) == g.cfg.GroupSize {
			g.incrementGroupCounter()
			currentGroup.GroupID = g.groupCounter

			go g.events.Handle(ctx, events.Message{Value: currentGroup})
		} else {
			remainingUsers = append(remainingUsers, currentGroup.Users...)
		}
	}

	if len(remainingUsers) > 0 {
		if err := g.service.SaveRemainingUsers(ctx, remainingUsers); err != nil {
			g.log.Error("error occurred on saving remaining users:", zap.String("error", err.Error()))
		}
	}
}

func (g *GroupFormationHandler) findClosestUsers(userMatrix []clusters.Coordinates, index int, groupedUsers map[int]struct{}) []dto.UserDistance {
	distances := make([]dto.UserDistance, len(userMatrix)-1)

	for j, userCoordinates := range userMatrix {
		_, exists := groupedUsers[j]
		if j == index || exists {
			continue
		}

		distance := g.euclideanDistance(userMatrix[index], userCoordinates)
		distances[j] = dto.UserDistance{Index: j, Distance: distance}
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].Distance < distances[j].Distance
	})

	return distances
}

func (g *GroupFormationHandler) formGroup(users []models.User, closestUsers []dto.UserDistance, groupedUsers map[int]struct{}) dto.Group {
	var group dto.Group
	group.Users = []models.User{}

	for _, closest := range closestUsers {
		if _, exists := groupedUsers[closest.Index]; !exists {
			groupedUsers[closest.Index] = struct{}{}
			group.Users = append(group.Users, users[closest.Index])

			if len(group.Users) == g.cfg.GroupSize {
				break
			}
		}
	}

	return group
}

func (g *GroupFormationHandler) euclideanDistance(a, b clusters.Coordinates) float64 {
	var sum float64
	for i := range a {
		d := a[i] - b[i]
		sum += d * d
	}
	return math.Sqrt(sum)
}

func (g *GroupFormationHandler) incrementGroupCounter() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.groupCounter++
}
