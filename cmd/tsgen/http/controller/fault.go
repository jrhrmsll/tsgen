package controller

import (
	"net/http"
	"strconv"

	"github.com/jrhrmsll/tsgen/pkg/model"
	"github.com/jrhrmsll/tsgen/pkg/store"

	"github.com/labstack/echo/v4"
)

const defaultIndent = "  "

type payload struct {
	Rate float32 `json:"rate"`
}

type Fault struct {
	store *store.Store
}

func NewFaultController(store *store.Store) *Fault {
	return &Fault{
		store: store,
	}
}

func (controller *Fault) Faults(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, controller.store.Faults(), defaultIndent)
}

func (controller *Fault) UpdateFault(c echo.Context) error {
	path := c.Param("path")

	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	payload := new(payload)
	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	fault, err := model.NewFault(path, code, payload.Rate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = controller.store.UpdateFault(fault)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
