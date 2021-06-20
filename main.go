package main

import (
	"os"

	"github.com/kaz/private-email-relay/internal/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err)
	}

	if err := s.Start(os.Getenv("K_SERVICE") == ""); err != nil {
		panic(err)
	}
}
