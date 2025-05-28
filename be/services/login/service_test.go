package login

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type ServiceTestSuite struct {
	suite.Suite
	ctx      context.Context
	mockRepo *MockUserRepo
	logger   *zap.Logger
	privKey  *rsa.PrivateKey
	service  *loginService
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctx = context.Background()
	ctrl := gomock.NewController(s.T())
	s.mockRepo = NewMockUserRepo(ctrl)
	s.logger = zap.NewNop()
	var err error
	s.privKey, err = rsa.GenerateKey(rand.Reader, 2048)
	s.Require().NoError(err)
	s.service = NewLoginService(s.mockRepo, s.privKey)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (s *ServiceTestSuite) TestLogin_Success() {
	req := RequestBody{Username: "user1", Password: "pass1"}
	s.mockRepo.EXPECT().
		FindByUsernameAndPassword(s.ctx, "user1", "pass1").
		Return("user-id-123", nil)

	token, err := s.service.Login(s.ctx, s.logger, req)
	s.NoError(err)
	s.NotEmpty(token)

}

func (s *ServiceTestSuite) TestLogin_UserNotFound() {
	req := RequestBody{Username: "user2", Password: "wrongpass"}
	s.mockRepo.EXPECT().
		FindByUsernameAndPassword(s.ctx, "user2", "wrongpass").
		Return("", errors.New("not found"))

	token, err := s.service.Login(s.ctx, s.logger, req)
	s.Error(err)
	s.Empty(token)
}

func (s *ServiceTestSuite) TestLogin_TokenSignError() {
	s.service.accessTokenPrivateKey = &rsa.PrivateKey{}

	req := RequestBody{Username: "user3", Password: "pass3"}
	s.mockRepo.EXPECT().
		FindByUsernameAndPassword(s.ctx, "user3", "pass3").
		Return("user-id-456", nil)

	token, err := s.service.Login(s.ctx, s.logger, req)
	s.Error(err)
	s.Empty(token)
}
