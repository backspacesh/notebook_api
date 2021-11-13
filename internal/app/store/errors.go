package store

import "errors"

var(
	ErrRecordNotFound = errors.New("record not found")
	ErrCreate = errors.New("create error")
)
