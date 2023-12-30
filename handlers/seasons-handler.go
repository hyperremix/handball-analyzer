package handlers

import (
	"github.com/hyperremix/handball-analyzer/components/pages"
	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/model"
	"github.com/hyperremix/handball-analyzer/services"
	"github.com/labstack/echo/v4"
)

func DefineSeasonsRoutes(e *echo.Echo) {
	e.GET("/", getSeasons)
}

func getSeasons(c echo.Context) error {
	var seasons []model.Season

	db.Get().Find(&seasons)

	return services.Render(c, pages.Home(seasons))
}
