package router

import "fmt"

type (
	Router interface {
		Set(from string, to string) error
		Unset(from string) error
	}
)

var (
	ErrorDuplicated       = fmt.Errorf("duplicated")
	ErrorUnsetNonexistent = fmt.Errorf("no such entry")
)
