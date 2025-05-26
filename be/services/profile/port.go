package profile

import (
	"context"

	"go.uber.org/zap"
)

type UserRepo interface {
	GetUserProfile(ctx context.Context, userID string) (string, error)
}

type Service interface {
	Profile(ctx context.Context, logger *zap.Logger, accessToken string) (string, error)
}
