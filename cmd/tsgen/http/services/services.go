package services

import (
	"math/rand"
	"net/http"

	"github.com/jrhrmsll/tsgen/pkg/model"
	"github.com/jrhrmsll/tsgen/pkg/store"

	"github.com/github/go-fault"
	"github.com/labstack/echo/v4"
)

const ALWAYS float32 = 1

type (
	Middleware  func(http.Handler) http.Handler
	Middlewares []Middleware
)

func (middlewares Middlewares) ToEchoMiddlewareFunc() []echo.MiddlewareFunc {
	echoMiddlewares := []echo.MiddlewareFunc{}
	for _, middleware := range middlewares {
		echoMiddlewares = append(echoMiddlewares, echo.WrapMiddleware(middleware))
	}

	return echoMiddlewares
}

type PathMiddlewareAdderService struct {
	store *store.Store
}

func NewPathMiddlewareAdderService(store *store.Store) *PathMiddlewareAdderService {
	return &PathMiddlewareAdderService{
		store: store,
	}
}

func (srv *PathMiddlewareAdderService) fn(fault model.Fault) func() float32 {
	// ALWAYS in this context means no participation
	return func() float32 {
		fault, err := srv.store.FindFaultBy(fault.Path, fault.Code)
		if err != nil {
			return ALWAYS
		}

		if rand.Float32() <= fault.Rate {
			return fault.Rate
		}

		return ALWAYS
	}
}

func (srv *PathMiddlewareAdderService) Adds(path model.Path) (Middlewares, error) {
	middlewares := Middlewares{}

	// slow injector is use to add some latency to the response
	slowInjector, err := fault.NewSlowInjector(path.ResponseTime)
	if err != nil {
		return nil, err
	}

	f, err := fault.NewFault(slowInjector,
		fault.WithEnabled(true),
		fault.WithParticipation(ALWAYS),
	)

	if err != nil {
		return nil, err
	}

	middlewares = append(middlewares, f.Handler)

	// path fauls are use to inject errors with a probability near the fault rate
	for _, pathFault := range path.Faults {
		errorInjector, err := fault.NewErrorInjector(pathFault.Code)
		if err != nil {
			return nil, err
		}

		f, err := fault.NewFault(errorInjector,
			fault.WithEnabled(true),
			fault.WithParticipation(ALWAYS),
			fault.WithRandFloat32Func(srv.fn(pathFault)),
		)

		if err != nil {
			return nil, err
		}

		middlewares = append(middlewares, f.Handler)
	}

	return middlewares, nil
}
