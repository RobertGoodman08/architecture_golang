package phoneNumber

import (
	"strings"
)

type PhoneNumber struct {
	value string
}

func (p PhoneNumber) String() string {
	return p.value
}

func New(phone string) *PhoneNumber {
	return &PhoneNumber{
		value: getNumbers(phone),
	}
}

func (p PhoneNumber) Equal(phoneNumber PhoneNumber) bool {
	return p.value == phoneNumber.value
}

func (p PhoneNumber) IsEmpty() bool {
	return len(strings.TrimSpace(p.value)) == 0
}
