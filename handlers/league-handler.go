package handlers

import (
	"context"

	"github.com/hyperremix/handball-analyzer/components/pages"
	"github.com/hyperremix/handball-analyzer/components/partials"
	"github.com/hyperremix/handball-analyzer/db"
	"github.com/hyperremix/handball-analyzer/environment"
	"github.com/hyperremix/handball-analyzer/services"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func DefineLeagueRoutes(e *echo.Echo) {
	e.GET("/seasons/:seasonID/leagues", getLeagues)
	e.GET("/seasons/:seasonID/leagues/:leagueID", getLeague)
}

type getBySeasonRequest struct {
	SeasonID int64 `param:"seasonID"`
}

func getLeagues(c echo.Context) error {
	var request getBySeasonRequest
	c.Bind(&request)
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, environment.DB_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	season, err := queries.GetSeason(ctx, request.SeasonID)
	if err != nil {
		return err
	}

	leagues, err := queries.ListSeasonLeagues(ctx, request.SeasonID)
	if err != nil {
		return err
	}

	leagueListResponse := services.MapToLeagueListResponse(season, leagues)

	isHxRequest := c.Request().Header.Get("HX-Request")

	if isHxRequest == "true" {
		return services.Render(c, partials.LeagueListPartial(leagueListResponse))
	}

	return services.Render(c, pages.LeagueListBase(leagueListResponse))
}

type getByLeagueRequest struct {
	LeagueID int64 `param:"leagueID"`
}

func getLeague(c echo.Context) error {
	var request getByLeagueRequest
	c.Bind(&request)

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, environment.DB_CONNECTION_STRING)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	league, err := queries.GetLeague(ctx, request.LeagueID)
	if err != nil {
		return err
	}

	games, err := queries.ListLeagueGames(ctx, league.ID)
	if err != nil {
		return err
	}

	teams, err := queries.ListLeagueTeams(ctx, league.ID)
	if err != nil {
		return err
	}

	isHxRequest := c.Request().Header.Get("HX-Request")

	details := services.MapToLeagueDetailsResponse(league, games, teams)

	if isHxRequest == "true" {
		return services.Render(c, partials.LeagueDetailsPartial(details))
	}

	return services.Render(c, pages.LeagueDetailsBase(details))
}
