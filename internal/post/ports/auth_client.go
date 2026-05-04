package ports

import (
	"context"
	"crypto/rsa"
)

type AuthClient interface {
	GetPublicKey(ctx context.Context) (*rsa.PublicKey, error)
}
