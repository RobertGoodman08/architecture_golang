package description

import "github.com/pkg/errors"

var (
	MaxLength      = 1000
	ErrWrongLength = errors.Errorf("description must be less than or equal to %d characters", MaxLength)
)

type Description struct {
	value string
}

func New(description string) (Description, error) {
	if len([]rune(description)) > MaxLength {
		return Description{}, ErrWrongLength
	}
	return Description{value: description}, nil
}

func (d Description) Value() string {
	return d.value
}
