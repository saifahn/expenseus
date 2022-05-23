package app

import "errors"

var (
	ErrDBItemNotFound = errors.New("store: item not found")
)
