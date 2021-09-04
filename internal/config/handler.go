package config

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Handler(cfg *Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, cfg.Raw)
	}
}
