package storage

import "errors"

var ErrNotFound = errors.New("storage: not found")
var ErrCodeExists = errors.New("storage: code exist")
