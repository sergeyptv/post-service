package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	User     string `env:"USER" env-prefix:"POSTGRES_" env-required`
	Password string `env:"PASSWORD" env-prefix:"POSTGRES_" env-required`
	Host     string `env:"HOST" env-prefix:"POSTGRES_" env-required`
	Port     string `env:"PORT" env-prefix:"POSTGRES_" env-required`
	DBName   string `env:"DBNAME" env-prefix:"POSTGRES_" env-required`
	// SslTls	string
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
	p.Close()
}
