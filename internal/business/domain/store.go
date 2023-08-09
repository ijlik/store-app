package domain

import (
	errpkg "github.com/ijlik/store-app/pkg/error"
	"time"
)

type Store struct {
	ID                   string    `json:"id"`
	Name                 string    `json:"name"`
	Url                  string    `json:"url"`
	Address              string    `json:"address"`
	Phone                string    `json:"phone"`
	OperationalTimeStart int       `json:"operational_time_start"`
	OperationalTimeEnd   int       `json:"operational_time_end"`
	CreatedAt            time.Time `json:"created_at"`
}

type StoreRequest struct {
	Name                 string `json:"name"`
	Url                  string `json:"-"`
	Address              string `json:"address"`
	Phone                string `json:"phone"`
	OperationalTimeStart int    `json:"operational_time_start"`
	OperationalTimeEnd   int    `json:"operational_time_end"`
}

func (s *StoreRequest) Validate() errpkg.ErrorService {
	if s.Name == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing name")
	}
	if s.Address == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing address")
	}
	if s.Phone == "" {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing phone")
	}
	if s.OperationalTimeStart > 23 || s.OperationalTimeStart < 0 {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing operational time start (0-23)")
	}
	if s.OperationalTimeEnd > 23 || s.OperationalTimeEnd < 0 {
		return errpkg.DefaultServiceError(errpkg.ErrBadRequest, "missing operational time end (0-23)")
	}
	s.Url = CreateSlug(s.Name, false)

	return nil
}

type HttpStoreIdParams struct {
	ID string `uri:"id"`
}
