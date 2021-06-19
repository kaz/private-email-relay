package main

import (
	"github.com/kaz/private-email-relay/internal/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err)
	}

	if err := s.Start(); err != nil {
		panic(err)
	}
}
