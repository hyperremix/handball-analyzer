package handlers

import "github.com/labstack/echo/v4"

func DefineLeagueRoutes(e *echo.Echo) {
	e.GET("/leagues", getLeagues)
}

func getLeagues(c echo.Context) error {
	return nil
}
