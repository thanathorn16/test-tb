package login

import (
	"context"

	"go.uber.org/zap"
)

//go:generate bash -c "mockgen -source=port.go -package=$(go list -f '{{.Name}}') -destination=port_mock_test.go"
type UserRepo interface {
	FindByUsernameAndPassword(ctx context.Context, username, password string) (string, error)
}

type Service interface {
	Login(ctx context.Context, logger *zap.Logger, req RequestBody) (string, error)
}
