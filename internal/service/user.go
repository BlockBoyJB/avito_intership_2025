package service

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/internal/repo"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
)

type userService struct {
	tx   *txmanager.TxManager
	user repo.User
}

func newUserService(tx *txmanager.TxManager, user repo.User) *userService {
	return &userService{
		tx:   tx,
		user: user,
	}
}

func (s *userService) SetIsActive(ctx context.Context, id string, status bool) (domain.User, error) {
	result, err := s.user.UpdateIsActive(ctx, id, status)
	return result, err
}
