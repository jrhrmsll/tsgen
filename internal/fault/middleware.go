package fault

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/github/go-fault"
	"github.com/labstack/echo/v4"
)

func Middleware(path string, code int, rate float32) (echo.MiddlewareFunc, error) {
	statusText := http.StatusText(code)
	if statusText == "" {
		return nil, fmt.Errorf("invalid fault code: %d", code)
	}

	var (
		rand = rand.New(rand.NewSource(upperBound))
		k    = key(path, code)
	)

	float32Func := func() float32 {
		v := rand.Float32()

		rate := store.get(k)
		if v <= rate {
			return rate
		}

		return upperBound
	}

	errorInjector, err := fault.NewErrorInjector(code)
	if err != nil {
		return nil, err
	}

	f, err := fault.NewFault(errorInjector,
		fault.WithEnabled(true),
		fault.WithParticipation(upperBound),
		fault.WithRandFloat32Func(float32Func),
	)

	if err != nil {
		return nil, err
	}

	store.set(k, rate)

	return echo.WrapMiddleware(f.Handler), nil
}
