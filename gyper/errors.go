package gyper

import (
	"errors"
)

var (
	ErrBinding = errors.New("failed to bind the request body to the provided struct")
)
