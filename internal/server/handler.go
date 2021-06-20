package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	RelayRequest struct {
		URL     string `json:"url"`
		Address string `json:"address"`
	}
)

func (s *Server) postRelay(c echo.Context) error {
	ctx := c.Request().Context()

	params := &RelayRequest{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
	}
	if params.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("`url` is required"))
	}

	addr, err := s.strategy.Assign(ctx, params.URL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to assign address: %v", err))
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
		"address": addr,
	})
}

func (s *Server) deleteRelay(c echo.Context) error {
	ctx := c.Request().Context()

	params := &RelayRequest{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
	}

	if params.URL != "" {
		if err := s.strategy.UnassignByUrl(ctx, params.URL); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign by url: %v", err))
		}
	} else if params.Address != "" {
		if err := s.strategy.UnassignByAddr(ctx, params.Address); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign by address: %v", err))
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("either `url` or `address` is required"))
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "deleted",
	})
}
