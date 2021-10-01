package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jrhrmsll/tsgen/cmd/tsgen/http/controller"
	"github.com/jrhrmsll/tsgen/cmd/tsgen/http/services"
	"github.com/jrhrmsll/tsgen/pkg/config"
	"github.com/jrhrmsll/tsgen/pkg/store"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const subsystem = "tsgen"

var configFile string

func main() {
	flag.StringVar(&configFile, "config.file", "config.yml", "tsgen configuration file path.")
	flag.Parse()

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	logger.Println("Starting tsgen Service")

	cfg, err := config.Load(configFile)
	if err != nil {
		logger.Fatal(err)
	}

	store, err := store.NewStore().Init(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	app := echo.New()
	{
		app.Use(middleware.Logger())
		app.Use(middleware.Recover())

		p := prometheus.NewPrometheus(subsystem, nil)
		p.Use(app)
	}

	// routes for paths and their faults as echo middlewares
	for _, path := range store.Paths() {
		middlewares, err := services.NewPathMiddlewareAdderService().Adds(path)
		if err != nil {
			logger.Fatal(err)
		}

		app.GET(path.Name, controller.Default(path.Name).Echo, middlewares.ToEchoMiddlewareFunc()...)
	}

	// API
	api := app.Group("api")
	{
		configController := controller.NewConfigController(cfg.Raw())
		faultController := controller.NewFaultController(store)

		api.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, "ok")
		})

		api.GET("/config", configController.Show)

		api.GET("/faults", faultController.Faults)
		api.PUT("/paths/:path/faults/:code", faultController.UpdateFault)
	}

	app.Logger.Fatal(app.Start(":8080"))
}
