package repo

import (
	"avito_intership_2025/internal/domain"
	"context"
)

type PullRequest interface {
	Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error)
	Find(ctx context.Context, id string) (domain.PullRequest, error)
	FindUserReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
	Merge(ctx context.Context, id string) error
	ReplaceReviewer(ctx context.Context, prID, oldReviewer, newReviewer string) error
	AddReviewers(ctx context.Context, prID string, reviewerID []string) error
	GetReviewers(ctx context.Context, prID string) ([]string, error)
}

type Team interface {
	Create(ctx context.Context, name string) error
	GetMembers(ctx context.Context, name string) ([]domain.User, error)
	GetUserMembers(ctx context.Context, userID string) ([]domain.User, error)
}

type User interface {
	Upsert(ctx context.Context, team string, users []domain.User) error
	UpdateIsActive(ctx context.Context, id string, status bool) (domain.User, error)
	Find(ctx context.Context, id string) (domain.User, error)
}

type Stats interface {
	GetReviewers(ctx context.Context) ([]domain.UserStats, error)
	GetPRByAuthor(ctx context.Context) ([]domain.AuthorStats, error)
	GetByTeam(ctx context.Context) ([]domain.TeamStats, error)
}

const (
	codeErrUniqueViolation = "23505"
)
