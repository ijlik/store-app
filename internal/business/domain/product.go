package domain

import (
	"fmt"
	errpkg "github.com/ijlik/store-app/pkg/error"
	"regexp"
	"strings"
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Price       float32   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Store       *Store    `json:"store"`
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Url         string  `json:"-"`
	Price       float32 `json:"price"`
	Description string  `json:"description"`
	StoreID     string  `json:"store_id"`
}

func (p *ProductRequest) Validate() errpkg.ErrorService {
	if p.Name == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing name")
	}
	if p.Price < 0 {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing price")
	}
	if p.Description == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing description")
	}
	if p.StoreID == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing store id")
	}
	p.Url = CreateSlug(p.Name, true)

	return nil
}

func CreateSlug(input string, isUnique bool) string {
	input = strings.ToLower(input)

	re := regexp.MustCompile("[^a-z0-9]+")
	input = re.ReplaceAllString(input, "-")
	input = strings.Trim(input, "-")

	if isUnique {
		unixCurrentTime := time.Now().Unix()
		return fmt.Sprintf("%s-%d", input, unixCurrentTime)
	} else {
		return input
	}
}

type HttpProductUrlParams struct {
	Url string `uri:"url"`
}

type HttpProductIdParams struct {
	ID string `uri:"id"`
}

type SearchAndFilterProduct struct {
	Limit         int    `form:"limit"`
	Page          int    `form:"page"`
	Search        string `form:"search"`
	SortDirection string `form:"sortDirection"`
	SortBy        string `form:"sortBy"`
}

func (sfe *SearchAndFilterProduct) Validate() errpkg.ErrorService {
	sfe.SortBy = getSortBy(sfe.SortBy)
	sfe.SortDirection = getSortDirection(sfe.SortDirection)

	if sfe.Limit <= 0 {
		sfe.Limit = 10
	}
	if sfe.Page <= 0 {
		sfe.Page = 1
	}

	return nil
}

var mapSortBy = map[string]string{
	"PRICE":      "price",
	"NAME":       "name",
	"CREATED_AT": "created_at",
}

func getSortBy(key string) string {
	item, ok := mapSortBy[strings.ToUpper(key)]
	if ok {
		return item
	}

	return mapSortBy["CREATED_AT"]
}

var mapSortDirection = map[string]string{
	"DESC": "DESC",
	"ASC":  "ASC",
}

func getSortDirection(key string) string {
	item, ok := mapSortDirection[strings.ToUpper(key)]
	if ok {
		return item
	}

	return mapSortDirection["DESC"]
}
