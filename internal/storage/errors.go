package storage

import "errors"

var (
	ErrNotFound = errors.New("note not found")
	ErrExists   = errors.New("note already exists")
)
