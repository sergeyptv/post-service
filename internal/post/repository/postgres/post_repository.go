package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/post/domain"
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

func (p *postgresPostRepository) Create(ctx context.Context, user domain.User, post domain.Post) (string, error) {
	const op = "repository.postgres.Create"

	var postUuid string

	err := p.db.QueryRow(ctx,
		`INSERT INTO post.article (user_uuid, username, description, media, created_at)
				VALUES ($1, $2, $3, $4, now())
				RETURNING uuid`,
		user.Uuid, user.Username, post.Description, post.Media).Scan(&postUuid)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return postUuid, nil
}

func (p *postgresPostRepository) Get(ctx context.Context, user domain.User, postUuid string) (domain.Post, error) {
	const op = "repository.postgres.Get"

	var post domain.Post

	err := p.db.QueryRow(ctx,
		`SELECT uuid, username, description, media, created_at, updated_at
				FROM post.article
				WHERE uuid = $1
					AND user_uuid = $2`,
		postUuid, user.Uuid).Scan(&post.Uuid, &post.Username, &post.Description, &post.Media, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return domain.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (p *postgresPostRepository) List(ctx context.Context, user domain.User) ([]string, error) {
	const op = "repository.postgres.List"

	var postUuid string
	postUuids := make([]string, 0)

	rows, err := p.db.Query(ctx,
		`SELECT uuid
				FROM post.article
				WHERE user_uuid = $1`,
		user.Uuid)
	if err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&postUuid)
		if err != nil {
			return []string{}, fmt.Errorf("%s: %w", op, err)
		}

		postUuids = append(postUuids, postUuid)
		postUuid = ""
	}

	err = rows.Err()
	if err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}

	return postUuids, nil
}

func (p *postgresPostRepository) Update(ctx context.Context, user domain.User, post domain.Post) error {
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

	vals = append(vals, post.Uuid, user.Uuid)

	cmdTag, err := p.db.Exec(ctx, query, vals)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() < 1 {
		return sql.ErrNoRows
	}

	return nil
}

func (p *postgresPostRepository) Delete(ctx context.Context, user domain.User, postUuid string) error {
	const op = "repository.postgres.Delete"

	cmdTag, err := p.db.Exec(ctx,
		`DELETE FROM post.article
				WHERE uuid = $1
					AND user_uuid = $2`,
		postUuid, user.Uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if cmdTag.RowsAffected() < 1 {
		return sql.ErrNoRows
	}

	return nil
}
