package login

import (
	"context"

	"go.uber.org/zap"
)

type UserRepo interface {
	FindByUsernameAndPassword(ctx context.Context, username, password string) (string, error)
	CreateUser(ctx context.Context, userID, username, password string) error
}

type Service interface {
	Login(ctx context.Context, logger *zap.Logger, req RequestBody) (string, error)
}
