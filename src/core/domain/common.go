package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type ID string

func NewID() ID {
	return ID(uuid.New().String())
}

// Errors

var (
	ErrValidation   = errors.New("[ValidationError]")
	ErrNotFound     = errors.New("[NotFound]")
	ErrNoPermission = errors.New("[NoPermission]")
	ErrAlreadyExist = errors.New("[AlreadyExist]")
	ErrUserNotFound = errors.New("[UserNotFound]")
)
