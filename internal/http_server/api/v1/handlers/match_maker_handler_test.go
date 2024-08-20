package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	DTOs "matchMaker/internal/dto"
	"matchMaker/internal/http_server/api/v1/handlers/mocks"
)

// mockery --name Service --dir internal/http_server/api/v1/handlers --output internal/http_server/api/v1/handlers/mocks --outpkg mocks

func TestSaveUsers(t *testing.T) {
	tests := []struct {
		name           string
		users          DTOs.Users
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Success",
			users:          getTestUsers(),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":201}`,
		},
		{
			name:           "Bad Request - Invalid Body",
			users:          DTOs.Users{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"error":"failed to validate user"}`,
		},
		{
			name:           "Internal Server Error",
			users:          getTestUsers(),
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"internal server error", "status":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.Service)
			if tt.mockError != nil {
				mockService.On("SaveUsers", mock.Anything, tt.users).Return(tt.mockError)
			} else {
				mockService.On("SaveUsers", mock.Anything, tt.users).Return(nil)
			}

			logger := zap.NewNop()
			v := validator.New()

			handler := NewMatchMakerHandler(logger, v, mockService)

			usersJSON, err := json.Marshal(tt.users)
			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(usersJSON))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Post("/api/v1/users", handler.SaveUsers(context.Background()))

			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.JSONEq(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
func getTestUsers() DTOs.Users {
	return DTOs.Users{
		User: []DTOs.User{
			{Name: "User 1", Skill: 5.1, Latency: 1.5},
			{Name: "User 2", Skill: 4.8, Latency: 7.1},
			{Name: "User 3", Skill: 3.9, Latency: 0.2},
		},
	}
}
