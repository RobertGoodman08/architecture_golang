package query

import (
	"strconv"
	"strings"

	"architecture_go/pkg/type/columnCode"
	"architecture_go/pkg/type/sort"
)

var (
	fieldSeparationCharacter = ","
	sortTypeCharacters       = []string{"-", "+"}

	defaultValueForLimit uint64 = 10
	maxValueForLimit     uint64 = 100
)

func parseSorts(strQuery string, options SortsOptions) (sort.Sorts, error) {
	var result = make(sort.Sorts, 0)
	if len(strQuery) == 0 {
		return result, nil
	}

	// Получаем массив из значений через запятую "-name,+full_name,phone" => ["-name", "+full_name", "phone"]
	for _, field := range strings.Split(strQuery, fieldSeparationCharacter) {

		if len(field) < 2 {
			continue
		}
		var name = field

		var direction = sort.DirectionAsc
		if strings.HasPrefix(field, sortTypeCharacters[1]) {
			direction = sort.DirectionAsc
			name = field[len(sortTypeCharacters[1]):]
		}

		if strings.HasPrefix(field, sortTypeCharacters[0]) {
			direction = sort.DirectionDesc
			name = field[len(sortTypeCharacters[0]):]
		}

		key, err := columnCode.New(name)
		if err != nil {
			return nil, err
		}

		if _, ok := options[key.String()]; !ok {
			continue
		}

		result = append(result, &sort.Sort{
			Key:       key,
			Direction: direction,
		})
	}

	return result, nil

}

func parseLimit(strLimit string) uint64 {
	limit, err := strconv.ParseUint(strLimit, 10, 64)
	if err != nil || limit == 0 {
		return defaultValueForLimit
	}

	if limit > maxValueForLimit {
		return maxValueForLimit
	}

	return limit
}

func parseOffset(strOffset string) uint64 {
	offset, err := strconv.ParseUint(strOffset, 10, 64)
	if err != nil {
		return 0
	}

	return offset
}
