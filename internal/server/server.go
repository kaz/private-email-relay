package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kaz/private-email-relay/internal/assign"
	"github.com/kaz/private-email-relay/internal/router"
	"github.com/kaz/private-email-relay/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Server struct {
		bindAddr string
		token    string

		assigners map[string]assign.Strategy
	}
)

func New() (*Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &Server{
		bindAddr: fmt.Sprintf(":%s", port),
	}

	server.token = os.Getenv("TOKEN")
	if server.token == "" {
		return nil, fmt.Errorf("TOKEN is missing")
	}

	var store storage.Storage
	if fsStore, err := storage.NewFirestoreStorage(context.Background()); err == nil {
		store = fsStore
	} else {
		fmt.Println("[[WARNING]] Using in-memory storage")
		store = storage.NewMemoryStorage()
	}

	var route router.Router
	if mgRoute, err := router.NewMailgunRouter(); err == nil {
		route = mgRoute
	} else {
		return nil, fmt.Errorf("no router is available: %w", err)
	}

	server.assigners = map[string]assign.Strategy{}
	if defaultAssign, err := assign.NewDefaultStrategy(store, route); err == nil {
		server.assigners["default"] = defaultAssign
	}
	if tempAssign, err := assign.NewTemporaryStrategy(store, route, func() time.Time { return time.Now().Add(3 * 24 * time.Hour) }); err == nil {
		server.assigners["temporary"] = tempAssign
	}
	if len(server.assigners) == 0 {
		return nil, fmt.Errorf("no strategy is available")
	}

	return server, nil
}

func (s *Server) Start(debug bool) error {
	e := echo.New()

	e.Debug = debug
	e.HidePort = !debug
	e.HideBanner = !debug

	e.Use(middleware.Logger())
	e.Use(s.authenticate)

	e.POST("/relay", s.postRelay)
	e.DELETE("/relay", s.deleteRelay)
	e.DELETE("/relay/expired", s.deleteRelayExpired)

	return e.Start(s.bindAddr)
}
