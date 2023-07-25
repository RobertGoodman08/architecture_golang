package sort

import "architecture_go/pkg/type/columnCode"

type Sort struct {
	Key columnCode.ColumnCode
	Direction
}

type Direction string

const (
	DirectionAsc  Direction = "ASC"
	DirectionDesc Direction = "DESC"
)

func (d Direction) String() string {
	return string(d)
}

func (s Sort) Parsing(mapping map[columnCode.ColumnCode]string) string {
	column, ok := mapping[s.Key]
	if !ok {
		return ""
	}
	return column + " " + s.Direction.String()
}

type Sorts []*Sort

func (s Sorts) Parsing(mapping map[columnCode.ColumnCode]string) []string {
	var result []string
	for _, sort := range s {

		result = append(result, sort.Parsing(mapping))
	}
	return result
}
