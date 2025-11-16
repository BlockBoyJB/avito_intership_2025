package service

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/repo"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
)

type PullRequest interface {
	Create(ctx context.Context, prID, name, authorID string) (domain.PullRequest, error)
	Merge(ctx context.Context, prID string) (domain.PullRequest, error)
	Reassign(ctx context.Context, prID, reviewerID string) (string, domain.PullRequest, error)
	GetUserReview(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

type Team interface {
	Create(ctx context.Context, team domain.Team) error
	Find(ctx context.Context, team string) (domain.Team, error)
}

type User interface {
	SetIsActive(ctx context.Context, id string, status bool) (domain.User, error)
}

type Stats interface {
	Reviewers(ctx context.Context) ([]domain.UserStats, error)
	ByAuthor(ctx context.Context) ([]domain.AuthorStats, error)
	ByTeam(ctx context.Context) ([]domain.TeamStats, error)
}

type Services struct {
	User        User
	Team        Team
	PullRequest PullRequest
	Stats       Stats
}

type ServicesDependencies struct {
	Tx          *txmanager.TxManager
	User        repo.User
	Team        repo.Team
	PullRequest repo.PullRequest
	Stats       repo.Stats
}

func NewServices(d *ServicesDependencies) *Services {
	return &Services{
		User:        newUserService(d.Tx, d.User),
		Team:        newTeamService(d.Tx, d.Team, d.User),
		PullRequest: newPRService(d.Tx, d.PullRequest, d.Team, d.User),
		Stats:       newStatsService(d.Stats),
	}
}
