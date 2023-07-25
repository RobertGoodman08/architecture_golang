package converter

import (
	"github.com/google/uuid"
)

func StringToUUID(str string) uuid.UUID {
	if len(str) == 0 {
		return uuid.Nil
	} else {
		if value, err := uuid.Parse(str); err != nil {
			return uuid.Nil
		} else {
			return value
		}
	}
}
