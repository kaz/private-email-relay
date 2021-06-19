package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getAddress(c echo.Context) error {
	addr, err := s.strategy.Assign(c.Request().Context(), c.QueryParam("url"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to assign address: %v", err))
	}
	return c.String(http.StatusCreated, addr)
}

func (s *Server) deleteAddress(c echo.Context) error {
	if err := s.strategy.Unassign(c.Request().Context(), c.QueryParam("addr")); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to unassign address: %v", err))
	}
	return c.NoContent(http.StatusNoContent)
}
