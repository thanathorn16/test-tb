package profile

import (
	"context"
	"crypto/rsa"

	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
)

type profileService struct {
	userRepo              UserRepo
	accessTokenPrivateKey *rsa.PrivateKey
}

func NewProfileService(userRepo UserRepo, accessTokenPrivateKey *rsa.PrivateKey) *profileService {
	return &profileService{userRepo: userRepo, accessTokenPrivateKey: accessTokenPrivateKey}
}

func (s *profileService) Profile(ctx context.Context, logger *zap.Logger, token string) (string, error) {
	// Parse the token
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return &s.accessTokenPrivateKey.PublicKey, nil
	})
	if err != nil {
		logger.Error("failed to parse token", zap.Error(err))
		return "", err
	}

	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if !ok || !tokenClaims.Valid {
		logger.Error("invalid token claims")
		return "", jwt.NewValidationError("invalid token claims", jwt.ValidationErrorClaimsInvalid)
	}

	userID, ok := claims["userID"].(string)
	if !ok {
		logger.Error("userID not found in token claims")
		return "", jwt.NewValidationError("userID not found in token claims", jwt.ValidationErrorClaimsInvalid)
	}

	// Fetch user profile using userID
	profile, err := s.userRepo.GetUserProfile(ctx, userID)
	if err != nil {
		logger.Error("failed to get user profile", zap.Error(err))
		return "", err
	}

	return profile, nil
}
