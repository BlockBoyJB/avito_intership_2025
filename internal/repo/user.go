package repo

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
)

type userRepo struct {
	tx txmanager.Manager
}

func NewUserRepo(tx txmanager.Manager) User {
	return &userRepo{tx: tx}
}

func (r *userRepo) Upsert(ctx context.Context, team string, users []domain.User) error {
	var (
		args   []any
		values []string
	)

	for i, u := range users {
		idx := i * 4
		values = append(values,
			fmt.Sprintf("($%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4),
		)
		args = append(args, u.ID, u.Username, team, u.IsActive)
	}

	sql := fmt.Sprintf(`
		INSERT INTO users (id, username, team_name, is_active)
		VALUES %s
		ON CONFLICT (id) DO UPDATE SET
			username = excluded.username,
			team_name = excluded.team_name,
			is_active = excluded.is_active;          
	`, strings.Join(values, ","))

	_, err := r.tx.Exec(ctx).Exec(ctx, sql, args...)
	return err
}

func (r *userRepo) UpdateIsActive(ctx context.Context, id string, status bool) (domain.User, error) {
	sql := "UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, username, team_name, is_active"

	var u domain.User
	if err := r.tx.Exec(ctx).QueryRow(ctx, sql, status, id).Scan(&u.ID, &u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
	}
	return u, nil
}

func (r *userRepo) Find(ctx context.Context, id string) (domain.User, error) {
	sql := "SELECT username, team_name, is_active FROM users WHERE id = $1"

	var u domain.User
	if err := r.tx.Exec(ctx).QueryRow(ctx, sql, id).Scan(&u.Username, &u.TeamName, &u.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}
