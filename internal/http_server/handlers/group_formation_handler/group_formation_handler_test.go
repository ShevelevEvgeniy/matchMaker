package group_formation_handler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"matchMaker/config"
	"matchMaker/internal/http_server/handlers/group_formation_handler/mocks"
	"matchMaker/internal/storage/postgres/repository/models"
)

//mockery --name Service --dir internal/http_server/handlers/group_formation_handler --output internal/http_server/handlers/group_formation_handler/mocks --outpkg mocks
//mockery --name Events --dir internal/http_server/handlers/group_formation_handler --output internal/http_server/handlers/group_formation_handler/mocks --outpkg mocks

func TestRun(t *testing.T) {
	mockService := new(mocks.Service)
	mockEvents := new(mocks.Events)

	mockService.On("GetUsersInSearch", mock.Anything, mock.Anything).Return([]models.User{{ID: 1}}, nil)
	mockService.On("GetAndRemoveRemainingUsers", mock.Anything).Return(nil, false, nil)

	handler := &GroupFormationHandler{
		cfg:     config.MatchSettings{CountWorkers: 1, BatchSize: 10, Delay: 100 * time.Millisecond},
		log:     zap.NewNop(),
		service: mockService,
		events:  mockEvents,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.Run(ctx)

	time.Sleep(200 * time.Millisecond)

	mockService.AssertExpectations(t)
}

func TestSelectUsers(t *testing.T) {
	mockService := new(mocks.Service)
	mockService.On("GetUsersInSearch", mock.Anything, mock.Anything).Return([]models.User{{ID: 1}}, nil)

	handler := &GroupFormationHandler{
		cfg:     config.MatchSettings{CountWorkers: 1, BatchSize: 10, Delay: 100 * time.Millisecond},
		log:     zap.NewNop(),
		service: mockService,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	usersChan := handler.selectUsers(ctx)

	users := <-usersChan
	assert.Len(t, users, 1)
	assert.Equal(t, int64(1), users[0].ID)
}

func TestGenerateGroups(t *testing.T) {
	mockService := new(mocks.Service)
	mockEvents := new(mocks.Events)

	mockService.On("GetUsersInSearch", mock.Anything, mock.Anything).Return([]models.User{{ID: 1}}, nil)
	mockService.On("GetAndRemoveRemainingUsers", mock.Anything).Return(nil, false, nil)
	mockService.On("SaveRemainingUsers", mock.Anything, mock.Anything).Return(nil)

	handler := &GroupFormationHandler{
		cfg:     config.MatchSettings{CountWorkers: 1, BatchSize: 10, GroupSize: 2},
		log:     zap.NewNop(),
		service: mockService,
		events:  mockEvents,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	usersChan := make(chan []models.User, 1)
	usersChan <- []models.User{{ID: 1}, {ID: 2}, {ID: 3}}

	close(usersChan)

	go handler.generateGroups(ctx, usersChan)

	time.Sleep(200 * time.Millisecond)

	mockService.AssertExpectations(t)
}
