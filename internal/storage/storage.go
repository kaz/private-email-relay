package storage

import "fmt"

type (
	Storage interface {
		Get(key string) (string, error)
		Set(key string, value string) error
		Unset(key string) error
	}
)

var (
	ErrorNotFound = fmt.Errorf("no such entry")
)
