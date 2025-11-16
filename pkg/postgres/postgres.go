package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"runtime"
	"time"
)

const (
	defaultConnAttempts = 10
	defaultConnTimeout  = 2 * time.Second
)

type postgres struct {
	*pgxpool.Pool
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
}

func NewPG(url string) (*pgxpool.Pool, error) {
	pg := &postgres{
		maxPoolSize:  runtime.NumCPU(),
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	poolConfig.MaxConns = int32(pg.maxPoolSize)
	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}
		log.Printf("Postgres trying to connect, attemps left: %d", pg.connAttempts)
		time.Sleep(pg.connTimeout)
		pg.connAttempts--
	}
	if err != nil {
		return nil, fmt.Errorf("error connect to postgres, %w", err)
	}
	return pg.Pool, err
}
