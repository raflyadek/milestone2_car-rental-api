package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)
func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := uuid.New().String()

		c.Request().Header.Set("X-request-id", requestId)

		log := logrus.WithFields(logrus.Fields{
			"request_id": requestId,
			"method": c.Request().Method,
			"path": c.Path(),
			"query": c.QueryString(),
			"remote_addr": c.Request().RemoteAddr,
			"user_agent": c.Request().UserAgent(),
			"referer": c.Request().Referer(),
		})

		start := time.Now()

		next(c)
		
		duration := time.Since(start)

		log.WithFields(logrus.Fields{
			"status_code": c.Request().Response.StatusCode,
			"duration_time": duration.Milliseconds(),
			"size_bytes": c.Response().Size,
		}).Info("HTTP request completed")

		return next(c)
	}
}