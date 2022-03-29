package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrDataAccess    = errors.New("data access")
	ErrSerialization = errors.New("serialization")
)

func NewNotFoundError(err interface{}) error {
	return fmt.Errorf("%w: %v", ErrNotFound, err)
}

func NewDataAccessError(err interface{}) error {
	return fmt.Errorf("%w: %v", ErrDataAccess, err)
}

func NewSerializationError(err interface{}) error {
	return fmt.Errorf("%w: %v", ErrSerialization, err)
}
