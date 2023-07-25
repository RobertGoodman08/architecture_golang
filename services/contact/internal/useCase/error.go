package useCase

import "errors"

var (
	ErrContactNotFound = errors.New("contact not found")
	ErrGroupNotFound   = errors.New("group not found")
)
