package txmanager

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txKey struct{}

type Querier interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
}

type Manager interface {
	Exec(ctx context.Context) Querier
}

type TxManager struct {
	pool *pgxpool.Pool
}

func NewManager(pool *pgxpool.Pool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

func (m *TxManager) Exec(ctx context.Context) Querier {
	tx, ok := m.getTx(ctx)
	if ok {
		return tx
	}
	return m.pool
}

func (m *TxManager) ExecInTx(ctx context.Context, f func(ctx context.Context) error) error {
	const op = "txmanager.ExecInTx"

	if tx, ok := m.getTx(ctx); tx != nil && ok {
		return f(ctx)
	}

	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s error begin tx: %w", op, err)
	}

	ctx = m.setTx(ctx, tx)
	var done bool
	defer func() {
		if !done {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = f(ctx); err != nil {
		return fmt.Errorf("%s error exec func: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s error commit tx: %w", op, err)
	}
	done = true
	return nil
}

func (m *TxManager) getTx(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(pgx.Tx)
	return tx, ok
}

func (m *TxManager) setTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}
