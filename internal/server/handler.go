package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/labstack/echo/v4"
)

type (
	PostRelayRequest struct {
		URL      string `json:"url"`
		Strategy string `json:"strategy"`
	}
	DeleteRelayRequst struct {
		URL      string `json:"url"`
		Address  string `json:"address"`
		Strategy string `json:"strategy"`
	}
)

func (s *Server) postRelay(c echo.Context) error {
	ctx := c.Request().Context()

	params := &PostRelayRequest{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
	}
	if params.URL == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("`url` is required"))
	}
	if params.Strategy == "" {
		params.Strategy = "default"
	}

	assigner, ok := s.assigners[params.Strategy]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no such strategy: %v", params.Strategy))
	}

	addr, err := assigner.Assign(ctx, params.URL)
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

	params := &DeleteRelayRequst{}
	if err := c.Bind(params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("failed to parse request: %v", err))
	}
	if params.Strategy == "" {
		params.Strategy = "default"
	}

	assigner, ok := s.assigners[params.Strategy]
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("no such strategy: %v", params.Strategy))
	}

	if params.URL != "" {
		if err := assigner.Unassign(ctx, params.URL); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign by url: %v", err))
		}
	} else if params.Address != "" {
		if err := assigner.UnassignByAddr(ctx, params.Address); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign by address: %v", err))
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("either `url` or `address` is required"))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (s *Server) deleteRelayExpired(c echo.Context) error {
	ctx := c.Request().Context()

	var tempAssigner *assign.TemporaryStrategy
	for _, assigner := range s.assigners {
		var ok bool
		if tempAssigner, ok = assigner.(*assign.TemporaryStrategy); ok {
			break
		}
	}

	count, err := tempAssigner.UnassignExpired(ctx, time.Now())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign expired address: %v", err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"count":   count,
	})
}
