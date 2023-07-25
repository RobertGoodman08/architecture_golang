package name

import (
	"github.com/pkg/errors"
)

var (
	MaxLength      = 250
	ErrWrongLength = errors.Errorf("name must be less than or equal to %d characters", MaxLength)
)

type Name struct {
	value string
}

func New(name string) (Name, error) {
	if len([]rune(name)) > MaxLength {
		return Name{}, ErrWrongLength
	}
	return Name{value: name}, nil
}

func (n Name) Value() string {
	return n.value
}
