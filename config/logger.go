package config

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func initializeLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func LoggerMiddleware() echo.MiddlewareFunc {
	l := initializeLogger()

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.Info().Str("Method", v.Method).Str("Uri", v.URI).Int("Status", v.Status).Msg("request")
			return nil
		},
	})
}
