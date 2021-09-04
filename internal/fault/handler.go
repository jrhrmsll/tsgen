package fault

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func List(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, store.faults(), defaultIndent)
}

func Update(c echo.Context) error {
	var (
		payload = new(struct {
			Rate float32 `json:"rate"`
		})
		path = c.Param("path")
	)

	code, err := strconv.Atoi(c.Param("code"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	statusText := http.StatusText(code)
	if statusText == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("invalid fault code: %d", code))
	}

	if err := c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	k := key(path, code)
	if store.has(k) {
		store.set(k, payload.Rate)
	}

	return c.NoContent(http.StatusNoContent)
}
