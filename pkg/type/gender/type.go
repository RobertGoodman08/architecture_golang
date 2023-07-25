package gender

type Gender uint8

func New(gender uint8) Gender {
	switch gender {
	case 1:
		return MALE
	case 2:
		return FEMALE
	default:
		return UNKNOWN
	}
}

const (
	UNKNOWN Gender = 0
	MALE    Gender = 1
	FEMALE  Gender = 2
)

func (g Gender) String() string {
	switch g {
	case 1:
		return "MALE"
	case 2:
		return "FEMALE"
	default:
		return "UNKNOWN"
	}
}

func (g Gender) Number() uint8 {
	return uint8(g)
}

func (g Gender) Equal(gender Gender) bool {
	return g == gender
}

func (g Gender) IsEmpty() bool {
	return g == UNKNOWN
}

func (g Gender) IsMale() bool {
	return g == MALE
}

func (g Gender) IsFemale() bool {
	return g == FEMALE
}
