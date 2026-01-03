package logger

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func EchoMiddleware(log *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			req := c.Request()
			res := c.Response()
			log.Info("http_request",
				slog.String("method", req.Method),
				slog.String("path", c.Path()),
				slog.Int("status", res.Status),
				slog.Duration("duration", time.Since(start)),
				slog.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			)
			return err
		}
	}
}
