package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sergeyptv/post_service/internal/auth/domain"
	"github.com/sergeyptv/post_service/internal/auth/repository"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
)

type postgresUserRepository struct {
	pool *postgres.Pool
}

func NewPostgresUserRepository(pool *postgres.Pool) *postgresUserRepository {
	return &postgresUserRepository{
		pool: pool,
	}
}

func (p *postgresUserRepository) CreateUser(ctx context.Context, user domain.CreateUser) (string, error) {
	const op = "repository.postgres.CreateUser"

	var userUuid string

	err := p.pool.Db.QueryRow(ctx,
		"INSERT INTO auth.users (username, passHash, email) VALUES ($1, $2, $3, $4) RETURNING uuid",
		user.Username, user.PassHash, user.Email,
	).Scan(&userUuid)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userUuid, nil
}

func (p *postgresUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "repository.postgres.GetUserByEmail"

	var user domain.User

	err := p.pool.Db.QueryRow(ctx,
		"SELECT uuid, username, passHash, email FROM auth.users WHERE email = $1",
		email).Scan(&user.Uuid, &user.Username, &user.PassHash, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}

		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
