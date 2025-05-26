package register

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"go.uber.org/zap"
)

type handler struct {
	service   Service
	logger    *zap.Logger
	validator *validator.Validate
}

func NewHandler(service Service, logger *zap.Logger) *handler {
	return &handler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

func (h *handler) Register(c echo.Context) error {
	var reqBody RequestBody

	if err := c.Bind(&reqBody); err != nil {
		h.logger.Error("request body binding failed", zap.Error(err))
		return c.JSON(http.StatusOK, MakeHTTPFailedResponse("9991", "invalid request"))
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		h.logger.Error("request body validation failed", zap.Error(err))
		return c.JSON(http.StatusOK, MakeHTTPFailedResponse("9991", "invalid request"))
	}

	err := h.service.Register(c.Request().Context(), h.logger, reqBody)
	if err != nil {
		h.logger.Error("register service failed", zap.String("error:", err.Error()))
		return c.JSON(http.StatusInternalServerError, MakeHTTPFailedResponse("9999", "internal server error"))
	}

	return c.JSON(http.StatusOK, MakeHTTPResponse("0000", "success", nil))

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
