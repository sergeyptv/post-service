package grpc

import (
	"context"
	"crypto/x509"
	authV1 "github.com/sergeyptv/post_service/api/pkg/proto/auth/v1"
	"github.com/sergeyptv/post_service/internal/auth/crypto/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	authV1.UnimplementedAuthServiceServer

	jwtConfig jwt.Config
}

func NewHandler(jwtConfig jwt.Config) *handler {
	return &handler{
		jwtConfig: jwtConfig,
	}
}

func (h *handler) GetPublicKey(ctx context.Context, req *authV1.GetPublicKeyRequest) (*authV1.GetPublicKeyResponse, error) {
	keyData, err := x509.MarshalPKIXPublicKey(h.jwtConfig.PublicKey)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create key data")
	}

	return &authV1.GetPublicKeyResponse{
		KeyData:   keyData,
		Format:    "DER",
		Algorithm: "RS256",
	}, nil
}
