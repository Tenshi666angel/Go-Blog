package db

import "errors"

const (
	UniqueConstraintCode = 1062
)

var (
	UniqueConstraintError = errors.New("unique constraint")
)
