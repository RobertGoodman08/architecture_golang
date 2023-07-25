package queryParameter

import (
	"architecture_go/pkg/type/pagination"
	"architecture_go/pkg/type/sort"
)

type QueryParameter struct {
	Sorts      sort.Sorts
	Pagination pagination.Pagination
	/*Тут можно добавить фильтр*/
}
