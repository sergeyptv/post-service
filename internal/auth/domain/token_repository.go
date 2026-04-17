package domain

import "context"

type TokenRepository interface {
	CreateToken(ctx context.Context, token Token) (string, error)
	GetToken(ctx context.Context, tokenUuid string) (Token, error)
	UpdateToken(ctx context.Context, tokenUuid string, updToken UpdateToken) error
	DeleteToken(ctx context.Context, tokenUuid string) error
}
