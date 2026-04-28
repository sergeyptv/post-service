package postgres

import (
	"context"
	"fmt"
	"github.com/sergeyptv/post_service/internal/platform/postgres"
	"github.com/sergeyptv/post_service/internal/post/domain"
)

type postgresPostRepository struct {
	db postgres.DBTX
}

func NewPostgresPostRepository(db postgres.DBTX) *postgresPostRepository {
	return &postgresPostRepository{
		db: db,
	}
}

func (p *postgresPostRepository) Create(ctx context.Context, post domain.Post) (string, error) {
	const op = "repository.postgres.Create"

	var postUuid string

	err := p.db.QueryRow(ctx,
		`INSERT INTO post.article (username, description, media, created_at)
				VALUES ($1, $2, $3, now())
				RETURNING uuid`,
		post.Username, post.Description, post.Media).Scan(&postUuid)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return postUuid, nil
}

func (p *postgresPostRepository) Get(ctx context.Context, postUuid string) (domain.Post, error) {
	const op = "repository.postgres.Get"

	var post domain.Post

	err := p.db.QueryRow(ctx,
		`SELECT uuid, username, description, media, created_at, updated_at
				FROM post.article
				WHERE uuid = $1`,
		postUuid).Scan(&post.Uuid, &post.Username, &post.Description, &post.Media, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return domain.Post{}, fmt.Errorf("%s: %w", op, err)
	}

	return post, nil
}

func (p *postgresPostRepository) List(ctx context.Context, username string) ([]string, error) {
	const op = "repository.postgres.List"

	var postUuid string
	postUuids := make([]string, 0)

	rows, err := p.db.Query(ctx,
		`SELECT uuid
				FROM post.article
				WHERE username = $1`,
		username)
	if err != nil {
		return []string{}, fmt.Errorf("%s: %w", op, err)
	}

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

func (p *postgresPostRepository) Update(ctx context.Context, post domain.Post) error {
	const op = "repository.postgres.Update"

	_, err := p.db.Exec(ctx,
		`UPDATE post.article
				SET description = $1, media = $2, updated_at = now()
				WHERE uuid = $3`,
		post.Description, post.Media, post.Uuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *postgresPostRepository) Delete(ctx context.Context, postUuid string) error {
	const op = "repository.postgres.Delete"

	_, err := p.db.Exec(ctx,
		`DELETE FROM post.article
				WHERE uuid = $1`,
		postUuid)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
