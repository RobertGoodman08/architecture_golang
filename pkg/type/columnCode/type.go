package columnCode

type ColumnCode string

func New(str string) (ColumnCode, error) {
	return ColumnCode(str), nil
}

func (c ColumnCode) String() string {
	return string(c)
}
