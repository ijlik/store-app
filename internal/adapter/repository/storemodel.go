package repository

import (
	"database/sql"
	"time"
)

type Store struct {
	ID                   string       `db:"id"`
	Name                 string       `db:"name"`
	Url                  string       `db:"url"`
	Address              string       `db:"address"`
	Phone                string       `db:"phone"`
	OperationalTimeStart int          `db:"operational_time_start"`
	OperationalTimeEnd   int          `db:"operational_time_end"`
	CreatedAt            time.Time    `db:"created_at"`
	UpdatedAt            sql.NullTime `db:"updated_at"`
}

func (s *Store) RowDataIndex() []interface{} {
	var data = []interface{}{
		s.ID,
		s.Name,
		s.Url,
		s.Address,
		s.Phone,
		s.OperationalTimeStart,
		s.OperationalTimeEnd,
		s.CreatedAt,
		s.UpdatedAt,
	}
	return data
}

func (s *Store) GetUpdatedAt() time.Time {
	return s.UpdatedAt.Time
}

func (s *Store) RowDataCreate() []interface{} {
	var data = []interface{}{
		s.Name,
		s.Url,
		s.Address,
		s.Phone,
		s.OperationalTimeStart,
		s.OperationalTimeEnd,
	}
	return data
}

func (s *Store) RowDataUpdate() []interface{} {
	var data = []interface{}{
		s.ID,
		s.Name,
		s.Url,
		s.Address,
		s.Phone,
		s.OperationalTimeStart,
		s.OperationalTimeEnd,
	}
	return data
}
