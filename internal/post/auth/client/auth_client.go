package client

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
)

var (
	ErrUnsupportedFormat    = errors.New("unsupported key format")
	ErrUnsupportedAlgorithm = errors.New("unsupported key sign algorithm")
	ErrParseRSAPublicKey    = errors.New("failed to parse RSA public key")
)

type authClient struct {
	client authV1.AuthServiceClient
}

func NewAuthClient(client authV1.AuthServiceClient) *authClient {
	return &authClient{
		client: client,
	}
}

func (a *authClient) GetPublicKey(ctx context.Context) (*rsa.PublicKey, error) {
	const op = "auth.client.GetPublicKey"

	resp, err := a.client.GetPublicKey(ctx, &authV1.GetPublicKeyRequest{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.GetAlgorithm() != "RS256" {
		return nil, fmt.Errorf("%s: %w", op, ErrUnsupportedAlgorithm)
	}

	switch resp.GetFormat() {
	case "DER":
		return a.parseDer(resp.GetKeyData())

	default:
		return nil, fmt.Errorf("%s: %w", op, ErrUnsupportedFormat)
	}
}

func (a *authClient) parseDer(keyData []byte) (*rsa.PublicKey, error) {
	const op = "auth.client.parseDer"

	pubInterface, err := x509.ParsePKIXPublicKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rsaPublicKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrParseRSAPublicKey)
	}

	return rsaPublicKey, nil
}
