package query

import (
	"github.com/gin-gonic/gin"

	"architecture_go/pkg/type/sort"
)

type Query struct {
	Sorts  sort.Sorts
	Limit  uint64
	Offset uint64
}

type SortOptions struct {
}

type Options struct {
	// Тут можно добавить фильтр
	Sorts SortsOptions
}

type SortsOptions map[string]SortOptions // map[front_key]FilterOptions

var (
	keyForSort        = "sort"
	defaultKeyForSort = ""
	keyForLimit       = "limit"
	keyForOffset      = "offset"
)

func ParseQuery(c *gin.Context, options Options) (*Query, error) {
	sorts, err := parseSorts(c.DefaultQuery(keyForSort, defaultKeyForSort), options.Sorts)
	if err != nil {
		return nil, err
	}

	return &Query{
		Sorts:  sorts,
		Limit:  parseLimit(c.Query(keyForLimit)),
		Offset: parseOffset(c.Query(keyForOffset)),
	}, nil
}

func ParseSorts(c *gin.Context, options SortsOptions) (sort.Sorts, error) {
	return parseSorts(c.DefaultQuery(keyForSort, defaultKeyForSort), options)
}

func ParseLimit(c *gin.Context) uint64 {
	return parseLimit(c.Query(keyForLimit))
}

func ParseOffset(c *gin.Context) uint64 {
	return parseOffset(c.Query(keyForOffset))
}
