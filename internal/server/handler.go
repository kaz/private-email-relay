package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getRelay(c echo.Context) error {
	ctx := c.Request().Context()

	url := c.QueryParam("url")
	if url == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("`url` is required"))
	}

	addr, err := s.strategy.Assign(ctx, url)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to assign address: %v", err))
	}
	return c.String(http.StatusCreated, addr)
}

func (s *Server) deleteRelay(c echo.Context) error {
	ctx := c.Request().Context()

	url := c.QueryParam("url")
	addr := c.QueryParam("addr")
	if url != "" {
		if err := s.strategy.UnassignByUrl(ctx, url); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign address: %v", err))
		}
	} else if addr != "" {
		if err := s.strategy.UnassignByAddr(ctx, addr); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign address: %v", err))
		}
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("either `url` or `addr` is required"))
	}
	return c.NoContent(http.StatusNoContent)
}
