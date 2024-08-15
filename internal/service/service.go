package service

import (
	"context"

	"matchMaker/internal/converter"
	"matchMaker/internal/dto"
	"matchMaker/internal/postgres/repository/models"
)

type Repository interface {
	Users(ctx context.Context, user models.User) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Users(ctx context.Context, user dto.User) error {
	return s.repo.Users(ctx, converter.ServiceToRepoModel(user))
}
