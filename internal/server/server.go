package server

import (
	"context"
	"fmt"
	"os"

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

		strategy assign.Strategy
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

	if dStrategy, err := assign.NewDefaultStrategy(store, route); err == nil {
		server.strategy = dStrategy
	} else {
		return nil, fmt.Errorf("no strategy is available: %w", err)
	}

	return server, nil
}

func (s *Server) Start() error {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(s.authenticate)

	e.GET("/address", s.getAddress)
	e.DELETE("/address", s.deleteAddress)

	return e.Start(s.bindAddr)
}
