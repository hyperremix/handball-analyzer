package handlers

import (
	"github.com/hyperremix/handball-analyzer/handballnet"
	"github.com/labstack/echo/v4"
)

func DefineTestRoutes(e *echo.Echo) {
	e.GET("/test", getTest)
}

func getTest(c echo.Context) error {
	err := handballnet.ProcessNewGames()

	if err != nil {
		return err
	}

	return c.JSON(200, nil)
}
