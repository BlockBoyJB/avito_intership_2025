package service

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/repo"
	"context"
)

type statsService struct {
	stats repo.Stats
}

func newStatsService(stats repo.Stats) *statsService {
	return &statsService{stats: stats}
}

func (s *statsService) Reviewers(ctx context.Context) ([]domain.UserStats, error) {
	return s.stats.GetReviewers(ctx)
}

func (s *statsService) ByAuthor(ctx context.Context) ([]domain.AuthorStats, error) {
	return s.stats.GetPRByAuthor(ctx)
}

func (s *statsService) ByTeam(ctx context.Context) ([]domain.TeamStats, error) {
	return s.stats.GetByTeam(ctx)
}
