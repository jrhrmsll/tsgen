package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"tsgen/internal/config"
	"tsgen/internal/fault"
	"tsgen/internal/slow"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const subsystem = "tsgen"

var configFile string

func success(v string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, v)
	}
}

func main() {
	flag.StringVar(&configFile, "config.file", "config.yml", "tsgen configuration file path.")
	flag.Parse()

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Starting tsgen Service")

	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	p := prometheus.NewPrometheus(subsystem, nil)
	p.Use(app)

	api := app.Group("api")
	{
		api.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "ok")
		})

		api.GET("/config", config.Handler(cfg))

		api.GET("/faults", fault.List)

		api.POST("/paths/:path/faults/:code", fault.Update)
	}

	for _, cfgPath := range cfg.Paths {
		middlewares := []echo.MiddlewareFunc{}
		for _, cfgFault := range cfgPath.Faults {
			errorMidleware, err := fault.Middleware(cfgPath.Name, cfgFault.Code, cfgFault.Rate)
			if err != nil {
				logger.Fatal(err)
			}

			middlewares = append(middlewares, errorMidleware)
		}

		slowMiddleware, err := slow.Middleware(cfgPath.ResponseTime)
		if err != nil {
			logger.Fatal(err)
		}

		middlewares = append(middlewares, slowMiddleware)

		app.GET(cfgPath.Name, success(cfgPath.Name), middlewares...)
	}

	app.Logger.Fatal(app.Start(":8080"))
}
