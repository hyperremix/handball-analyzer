package handlers

import (
	"context"

	"github.com/hyperremix/handball-analyzer/components/pages"
	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/environment"
	"github.com/hyperremix/handball-analyzer/services"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func DefineSeasonsRoutes(e *echo.Echo) {
	e.GET("/", getSeasons)
}

func getSeasons(c echo.Context) error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, environment.DB_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	seasonRows, err := queries.ListSeasons(ctx)
	if err != nil {
		return err
	}

	seasonResponses := services.MapRowsToSeasonResponses(seasonRows)

	return services.Render(c, pages.SeasonsBase(seasonResponses))
}
