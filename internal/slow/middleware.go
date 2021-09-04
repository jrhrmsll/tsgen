package slow

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	upperBound   = 100
	milliseconds = 1_000_000 // 1ms = 1,000,000ns
)

func Middleware(t time.Duration) (echo.MiddlewareFunc, error) {
	rand := rand.New(rand.NewSource(time.Now().Unix()))

	handlerFunc := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			v := rand.Intn(upperBound) * milliseconds
			time.Sleep(t + time.Duration(v))

			return next(c)
		}
	}

	return handlerFunc, nil
}
