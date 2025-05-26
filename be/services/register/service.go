package register

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type registerService struct {
	userRepo UserRepo
}

func NewRegisterService(userRepo UserRepo) *registerService {
	return &registerService{userRepo: userRepo}
}

func (s *registerService) Register(ctx context.Context, logger *zap.Logger, req RequestBody) error {
	exists, err := s.userRepo.CheckUserExists(ctx, req.Username)
	if err != nil {
		logger.Error("failed to check if user exists", zap.Error(err))
		return err
	}
	if exists {
		logger.Error("user already exists", zap.String("username", req.Username))
		return errors.New("user already exists")
	}

	userID := uuid.New().String()
	err = s.userRepo.CreateUser(ctx, userID, req.Username, req.Password)
	if err != nil {
		logger.Error("failed to create user", zap.Error(err))
		return err
	}
	return nil
}
