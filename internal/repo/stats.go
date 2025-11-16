package repo

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
	"github.com/jackc/pgx/v5"
)

type statsRepo struct {
	tx txmanager.Manager
}

func NewStatsRepo(tx txmanager.Manager) Stats {
	return &statsRepo{tx: tx}
}

func (r *statsRepo) GetReviewers(ctx context.Context) ([]domain.UserStats, error) {
	sql := `
		SELECT reviewer_id, COUNT(*) AS assignments
		FROM pr_reviewers
		GROUP BY reviewer_id
		ORDER BY assignments DESC
	`
	rows, err := r.tx.Exec(ctx).Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.UserStats])
}

func (r *statsRepo) GetPRByAuthor(ctx context.Context) ([]domain.AuthorStats, error) {
	sql := `
		SELECT author_id, COUNT(*) AS pr_count 
		FROM pull_requests
		GROUP BY author_id
		ORDER BY pr_count DESC	
	`
	rows, err := r.tx.Exec(ctx).Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.AuthorStats])
}

func (r *statsRepo) GetByTeam(ctx context.Context) ([]domain.TeamStats, error) {
	sql := `
		SELECT team_name, COUNT(*) AS members 
		FROM users
		GROUP BY team_name
		ORDER BY members DESC	
	`
	rows, err := r.tx.Exec(ctx).Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.TeamStats])
}
