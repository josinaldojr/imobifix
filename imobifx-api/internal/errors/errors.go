package errors

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string { return e.Code + ": " + e.Message }

func New(status int, code, msg string, details interface{}) *AppError {
	return &AppError{Status: status, Code: code, Message: msg, Details: details}
}

func FiberErrorHandler(c *fiber.Ctx, err error) error {
	rid, _ := c.Locals("requestid").(string)

	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Status >= 500 {
			slog.Error("app_error",
				slog.String("request_id", rid),
				slog.String("code", appErr.Code),
				slog.Int("status", appErr.Status),
				slog.String("path", c.Path()),
				slog.Any("details", appErr.Details),
			)
		}
		return c.Status(appErr.Status).JSON(appErr)
	}

	var fe *fiber.Error
	if errors.As(err, &fe) {
		return c.Status(fe.Code).JSON(New(fe.Code, "HTTP_ERROR", fe.Message, nil))
	}

	slog.Error("unhandled_error",
		slog.String("request_id", rid),
		slog.String("path", c.Path()),
		slog.String("method", c.Method()),
		slog.String("error", err.Error()),
	)

	return c.Status(http.StatusInternalServerError).
		JSON(New(http.StatusInternalServerError, "INTERNAL_ERROR", "Erro interno.", nil))
}
