package service

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/repo"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
)

type teamService struct {
	tx   *txmanager.TxManager
	team repo.Team
	user repo.User
}

func newTeamService(tx *txmanager.TxManager, team repo.Team, user repo.User) *teamService {
	return &teamService{
		tx:   tx,
		team: team,
		user: user,
	}
}

func (s *teamService) Create(ctx context.Context, team domain.Team) error {
	err := s.tx.ExecInTx(ctx, func(ctx context.Context) error {
		if err := s.team.Create(ctx, team.Name); err != nil {
			return err
		}
		if err := s.user.Upsert(ctx, team.Name, team.Members); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *teamService) Find(ctx context.Context, team string) (domain.Team, error) {
	users, err := s.team.GetMembers(ctx, team)
	if err != nil {
		return domain.Team{}, err
	}
	return domain.Team{
		Name:    team,
		Members: users,
	}, nil
}
