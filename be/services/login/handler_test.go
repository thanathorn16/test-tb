package login

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"
	gomock "go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type HandlerTestSuite struct {
	suite.Suite
	ctx context.Context
	e   *echo.Echo

	handler  *handler
	mService *MockService
	mLogger  *zap.Logger
}

func TestBackgroundTestSuiteHandler(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (s *HandlerTestSuite) SetupTest() {
	s.ctx = context.Background()
	mCtrl := gomock.NewController(s.T())
	s.mService = NewMockService(mCtrl)

	s.mLogger = zap.NewNop()
	s.handler = NewHandler(s.mService, s.mLogger)

	s.e = echo.New()
}

func (s *HandlerTestSuite) TestLogin_Success() {
	body := `{"username":"user1","password":"pass1"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	s.mService.EXPECT().
		Login(gomock.Any(), gomock.Any(), RequestBody{Username: "user1", Password: "pass1"}).
		Return("token123", nil)

	err := s.handler.Login(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), `"code":"0000"`)
	s.Contains(rec.Body.String(), `"token":"token123"`)
}

func (s *HandlerTestSuite) TestLogin_BindError() {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{invalid json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	err := s.handler.Login(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), `"code":"9991"`)
}

func (s *HandlerTestSuite) TestLogin_ValidationError() {
	body := `{"username":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	err := s.handler.Login(c)
	s.NoError(err)
	s.Equal(http.StatusOK, rec.Code)
	s.Contains(rec.Body.String(), `"code":"9991"`)
}

func (s *HandlerTestSuite) TestLogin_ServiceError() {
	body := `{"username":"user1","password":"pass1"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := s.e.NewContext(req, rec)

	s.mService.EXPECT().
		Login(gomock.Any(), gomock.Any(), RequestBody{Username: "user1", Password: "pass1"}).
		Return("", errors.New("service error"))

	err := s.handler.Login(c)
	s.NoError(err)
	s.Equal(http.StatusInternalServerError, rec.Code)
	s.Contains(rec.Body.String(), `"code":"9999"`)
}
