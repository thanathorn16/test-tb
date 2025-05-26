package profile

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *handler {
	return &handler{
		service: service,
		logger:  logger,
	}
}

func (h *handler) Register(c echo.Context) error {

	authToken := c.Request().Header.Get("Authorization")
	token := strings.TrimPrefix(authToken, "Bearer ")
	if token == "" {
		h.logger.Error("authorization token is missing")
		return c.JSON(http.StatusUnauthorized, MakeHTTPFailedResponse("9991", "invalid request"))
	}

	resp, err := h.service.Profile(c.Request().Context(), h.logger, token)
	if err != nil {
		h.logger.Error("get profile user failed", zap.String("error:", err.Error()))
		return c.JSON(http.StatusInternalServerError, MakeHTTPFailedResponse("9999", "Internal Server Error"))
	}

	return c.JSON(http.StatusOK, MakeHTTPResponse("0000", "success", map[string]string{"userName": resp}))

}

func MakeHTTPResponse(code, message string, data interface{}) *httpResponseTemplate {
	return &httpResponseTemplate{
		Result: Result{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
}

func MakeHTTPFailedResponse(code, message string) *httpResponseTemplate {
	return &httpResponseTemplate{
		Result: Result{
			Code:    code,
			Message: message,
		},
	}
}
