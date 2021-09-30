package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Default string

func NewDefaultController(v string) Default {
	return Default(v)
}

func (controller Default) Echo(c echo.Context) error {
	return c.String(http.StatusOK, string(controller))
}
