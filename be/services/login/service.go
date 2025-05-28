package login

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type loginService struct {
	userRepo              UserRepo
	accessTokenPrivateKey *rsa.PrivateKey
}

func NewLoginService(userRepo UserRepo, accessTokenPrivateKey *rsa.PrivateKey) *loginService {
	return &loginService{userRepo: userRepo, accessTokenPrivateKey: accessTokenPrivateKey}
}

func (s *loginService) Login(ctx context.Context, logger *zap.Logger, req RequestBody) (string, error) {
	userID, err := s.userRepo.FindByUsernameAndPassword(ctx, req.Username, req.Password)
	if err != nil {
		logger.Error("failed to find user by username and password", zap.Error(err))
		return "", err
	}

	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 1).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.accessTokenPrivateKey)
	if err != nil {
		logger.Error("failed to sign token", zap.Error(err))
		return "", err
	}

	return tokenString, nil
}
