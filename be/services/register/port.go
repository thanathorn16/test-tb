package register

import (
	"context"

	"go.uber.org/zap"
)

type UserRepo interface {
	CheckUserExists(ctx context.Context, username string) (bool, error)
	CreateUser(ctx context.Context, userID, username, password string) error
}

type Service interface {
	Register(ctx context.Context, logger *zap.Logger, req RequestBody) error
}
