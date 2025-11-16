package repo

import (
	"avito_intership_2025/internal/domain"
	"avito_intership_2025/pkg/postgres/txmanager"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type prRepo struct {
	tx txmanager.Manager
}

func NewPRRepo(tx txmanager.Manager) PullRequest {
	return &prRepo{tx: tx}
}

func (r *prRepo) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	sql := "INSERT INTO pull_requests (id, name, author_id) VALUES ($1, $2, $3) RETURNING id, name, author_id, status, created_at, merged_at"

	var out domain.PullRequest
	err := r.tx.Exec(ctx).QueryRow(ctx, sql, pr.ID, pr.Name, pr.AuthorID).Scan(
		&out.ID,
		&out.Name,
		&out.AuthorID,
		&out.Status,
		&out.CreatedAt,
		&out.MergedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == codeErrUniqueViolation {
				return domain.PullRequest{}, domain.ErrPRExists
			}
		}
	}
	return out, nil
}

func (r *prRepo) Find(ctx context.Context, id string) (domain.PullRequest, error) {
	sql := `
		SELECT pr.id, name, author_id, status, created_at, merged_at,
		(
			SELECT COALESCE(array_agg(reviewer_id), '{}')
			FROM pr_reviewers r
			WHERE r.pr_id = pr.id
		) AS reviewers
		FROM pull_requests pr
		WHERE pr.id = $1
	`

	rows, err := r.tx.Exec(ctx).Query(ctx, sql, id)
	if err != nil {
		return domain.PullRequest{}, err
	}
	pr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.PullRequest])
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.PullRequest{}, domain.ErrPRNotFound
	}
	return pr, nil
}

func (r *prRepo) FindUserReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	sql := `
		SELECT pr.id, name, author_id, status , created_at, merged_at
		FROM pr_reviewers r
		JOIN pull_requests pr on r.pr_id = pr.id
		WHERE r.reviewer_id = $1
		ORDER BY created_at
	`
	rows, err := r.tx.Exec(ctx).Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[domain.PullRequest])
}

func (r *prRepo) Merge(ctx context.Context, id string) error {
	sql := "UPDATE pull_requests SET status = 'MERGED', merged_at = COALESCE(merged_at, now()) WHERE id = $1 RETURNING id, name, author_id, status, merged_at"

	result, err := r.tx.Exec(ctx).Exec(ctx, sql, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrPRNotFound
	}
	return nil
}

func (r *prRepo) ReplaceReviewer(ctx context.Context, prID, oldReviewer, newReviewer string) error {
	sql := "UPDATE pr_reviewers SET reviewer_id = $1 WHERE pr_id = $2 AND reviewer_id = $3"
	_, err := r.tx.Exec(ctx).Exec(ctx, sql, newReviewer, prID, oldReviewer)
	return err
}

func (r *prRepo) AddReviewers(ctx context.Context, prID string, reviewerID []string) error {
	sql := "INSERT INTO pr_reviewers (pr_id, reviewer_id) SELECT $1, unnest($2::varchar[])"

	if _, err := r.tx.Exec(ctx).Exec(ctx, sql, prID, reviewerID); err != nil {
		return err
	}
	return nil
}

func (r *prRepo) GetReviewers(ctx context.Context, prID string) ([]string, error) {
	sql := "SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1"

	rows, err := r.tx.Exec(ctx).Query(ctx, sql, prID)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowTo[string])
}
