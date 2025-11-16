package repo

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type teamRepo struct {
	tx txmanager.Manager
}

func NewTeamRepo(tx txmanager.Manager) Team {
	return &teamRepo{tx: tx}
}

func (r *teamRepo) Create(ctx context.Context, name string) error {
	sql := "INSERT INTO teams (name) VALUES ($1)"
	if _, err := r.tx.Exec(ctx).Exec(ctx, sql, name); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == codeErrUniqueViolation {
				return domain.ErrTeamExists
			}
		}
		return err
	}
	return nil
}

func (r *teamRepo) GetMembers(ctx context.Context, name string) ([]domain.User, error) {
	sql := "SELECT id, username, is_active FROM users WHERE team_name = $1"

	return r.collectUsers(ctx, sql, name)
}

func (r *teamRepo) GetUserMembers(ctx context.Context, userID string) ([]domain.User, error) {
	sql := "SELECT id, username, is_active FROM users WHERE team_name = (SELECT team_name FROM users WHERE id = $1)"

	return r.collectUsers(ctx, sql, userID)
}

func (r *teamRepo) collectUsers(ctx context.Context, sql string, args ...any) ([]domain.User, error) {
	rows, err := r.tx.Exec(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.User])
}
