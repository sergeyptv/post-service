package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/sergeyptv/post_service/platform/postgres"
	"github.com/sergeyptv/post_service/post/internal/domain"
	"strings"
)

type postgresPostRepository struct {
	db postgres.DBTX
}

func NewPostgresPostRepository(db postgres.DBTX) *postgresPostRepository {
	return &postgresPostRepository{
		db: db,
	}
}

func (p *postgresPostRepository) Create(ctx context.Context, userUuid, username string, post domain.Post) (string, error) {
	const op = "repository.postgres.Create"

	var postUuid string

	err := p.db.QueryRow(ctx,
		`INSERT INTO post.article (user_uuid, username, description, media, created_at)
				VALUES ($1, $2, $3, $4, now())
				RETURNING uuid`,
		userUuid, username, post.Description, post.Media).Scan(&postUuid)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return postUuid, nil
}

func (p *postgresPostRepository) Get(ctx context.Context, userUuid string, postUuid string) (domain.Post, error) {
	const op = "repository.postgres.Get"

	var post domain.Post

	err := p.db.QueryRow(ctx,
		`SELECT uuid, username, description, media, created_at, updated_at
				FROM post.article
				WHERE uuid = $1
					AND user_uuid = $2`,
		postUuid, userUuid).Scan(&post.Uuid, &post.Username, &post.Description, &post.Media, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Post{}, domain.ErrPostNotExist
		}

		return domain.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (p *postgresPostRepository) List(ctx context.Context, userUuid string) ([]string, error) {
	const op = "repository.postgres.List"

	postUuids := make([]string, 0)

	rows, err := p.db.Query(ctx,
		`SELECT uuid
				FROM post.article
				WHERE user_uuid = $1`,
		userUuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, domain.ErrPostNotExist
		}

		return []string{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var postUuid string

		err = rows.Scan(&postUuid)
		if err != nil {
			return []string{}, fmt.Errorf("%s: %w", op, err)
		}

		postUuids = append(postUuids, postUuid)
	}

	err = rows.Err()
	if err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}

	return postUuids, nil
}

func (p *postgresPostRepository) Update(ctx context.Context, userUuid string, post domain.Post) error {
	const op = "repository.postgres.Update"

	setParts := make([]string, 0)
	vals := make([]any, 0)
	idx := 1

	if post.Description != "" {
		setParts = append(setParts, fmt.Sprintf("description = $%d", idx))
		vals = append(vals, post.Description)
		idx++
	}
	if len(post.Media) > 0 {
		setParts = append(setParts, fmt.Sprintf("media = $%d", idx))
		vals = append(vals, post.Media)
		idx++
	}

	query := "UPDATE post.article SET " + strings.Join(setParts, ", ")
	query += fmt.Sprintf(" WHERE uuid = $%d AND user_uuid = $%d", idx, idx+1)

	vals = append(vals, post.Uuid, userUuid)

	cmdTag, err := p.db.Exec(ctx, query, vals)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() < 1 {
		return pgx.ErrNoRows
	}

	return nil
}

func (p *postgresPostRepository) Delete(ctx context.Context, userUuid, postUuid string) error {
	const op = "repository.postgres.Delete"

	cmdTag, err := p.db.Exec(ctx,
		`DELETE FROM post.article
				WHERE uuid = $1
					AND user_uuid = $2`,
		postUuid, userUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() < 1 {
		return domain.ErrPostNotExist
	}

	return nil
}
