package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Config struct {
	raw string
}

func NewConfigController(raw string) *Config {
	return &Config{
		raw: raw,
	}
}

func (controller *Config) Show(c echo.Context) error {
	return c.String(http.StatusOK, controller.raw)
}
