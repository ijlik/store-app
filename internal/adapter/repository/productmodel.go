package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Product struct {
	ID          string       `db:"id"`
	StoreID     string       `db:"store_id"`
	Name        string       `db:"name"`
	Url         string       `db:"url"`
	Price       float32      `db:"price"`
	Description string       `db:"description"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

func (p *Product) RowDataIndex() []interface{} {
	var data = []interface{}{
		p.ID,
		p.StoreID,
		p.Name,
		p.Url,
		p.Price,
		p.Description,
		p.CreatedAt,
		p.UpdatedAt,
	}
	return data
}

func (p *Product) GetUpdatedAt() time.Time {
	return p.UpdatedAt.Time
}

func (p *Product) RowDataCreate() []interface{} {
	var data = []interface{}{
		p.StoreID,
		p.Name,
		p.Url,
		p.Price,
		p.Description,
	}
	return data
}

func (p *Product) RowDataUpdate() []interface{} {
	var data = []interface{}{
		p.ID,
		p.StoreID,
		p.Name,
		p.Url,
		p.Price,
		p.Description,
	}
	return data
}

type SearchFilterPagination struct {
	Limit         int
	Offset        int
	Search        string
	SortDirection string
	SortBy        string
}

func (sfp *SearchFilterPagination) BuildWhere(baseQuery string, usePagination bool, customCondition string) (string, []any, error) {
	var (
		query       = baseQuery + " WHERE 1=1"
		params      []interface{}
		defSearchBy = []string{
			"name",
		}
		paramIndex = 1
	)
	if customCondition != "" {
		query = query + " AND " + customCondition
	}

	if sfp.Search != "" {
		query = query + " AND ("
		for _, field := range defSearchBy {
			query = query + fmt.Sprintf("%s ILIKE $%d OR ", field, paramIndex)
			params = append(params, "%"+sfp.Search+"%")
			paramIndex++
		}
		query = strings.TrimRight(query, " OR ") + ")"
	}

	if usePagination && sfp.SortBy != "" {
		query = query + fmt.Sprintf(" ORDER BY %s %s", sfp.SortBy, sfp.SortDirection)
	}
	if usePagination && sfp.Limit != 0 {
		query = query + fmt.Sprintf(" LIMIT %d OFFSET %d", sfp.Limit, sfp.Offset)
	}

	return query, params, nil
}
