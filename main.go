package main

import (
	"github.com/hyperremix/handball-analyzer/components/pages"
	"github.com/hyperremix/handball-analyzer/components/partials"
	"github.com/hyperremix/handball-analyzer/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status}\t| ${latency_human}\t| ${method} | ${uri}\n",
	}))
	e.Static("/assets", "./assets")
	e.GET("/", func(c echo.Context) error {
		return services.Render(c, pages.Home())
	})

	e.GET("/seasons", func(c echo.Context) error {
		return services.Render(c, partials.Seasons())
	})

	e.Logger.Fatal(e.Start("localhost:8080"))
}
