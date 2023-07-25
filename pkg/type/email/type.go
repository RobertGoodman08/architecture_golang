package email

import (
	"errors"
	"regexp"
	"strings"
)

var regexpEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

type Email struct {
	value string
}

func New(email string) (Email, error) {
	if len(email) == 0 {
		return Email{}, nil
	}

	if !regexpEmail.MatchString(strings.ToLower(email)) {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: email}, nil
}

func (e Email) Email() Email {
	return e
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equal(email Email) bool {
	return e.value == email.value
}

func (e Email) IsEmpty() bool {
	return len(strings.TrimSpace(e.String())) == 0
}

func (e *Email) MarshalJSON() ([]byte, error) {
	return []byte(`"` + e.value + `"`), nil
}

func (e *Email) UnmarshalJSON(bin []byte) error {
	str := string(bin)
	if strings.HasPrefix(str, `"`) {
		str = str[1:]
	}
	if strings.HasSuffix(str, `"`) {
		str = str[:len(str)-1]
	}
	tmp, err := New(str)
	if err != nil {
		return err
	}
	e.value = tmp.value
	return nil
}
