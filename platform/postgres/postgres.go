package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `env:"USER" env-required`
	Password string `env:"PASSWORD" env-required`
	Host     string `env:"HOST" env-required`
	Port     string `env:"PORT" env-required`
	DBName   string `env:"DBNAME" env-required`
}

type DBTX interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Pool struct {
	Db *pgxpool.Pool
}

func NewPool(ctx context.Context, c Config) (*Pool, error) {
	pool, err := pgxpool.New(
		ctx,
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", c.User, c.Password, c.Host, c.Port, c.DBName),
	)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &Pool{
		Db: pool,
	}, nil
}

func (p *Pool) Close() {
	p.Db.Close()
}
