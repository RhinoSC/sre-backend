package internal

import "errors"

type Prize struct {
	ID                    string
	Name                  string
	Description           string
	Url                   string
	MinAmount             float64
	Status                string
	InternationalDelivery bool
	EventID               string
}

var (
	// REPOSITORY ERRORS
	ErrPrizeRepositoryNotFound   = errors.New("repository: prize not found")
	ErrPrizeRepositoryDuplicated = errors.New("repository: prize already exists")
	// SERVICE ERRORS
	ErrPrizeDatabase          = errors.New("database error")
	ErrPrizeServiceNotFound   = errors.New("service: prize not found")
	ErrPrizeServiceDuplicated = errors.New("service: prize already exists")
)

type PrizeRepository interface {
	FindAll() ([]Prize, error)

	FindById(id string) (Prize, error)

	Save(prize *Prize) error

	Update(prize *Prize) error

	Delete(id string) error
}

type PrizeService interface {
	FindAll() ([]Prize, error)

	FindById(id string) (Prize, error)

	Save(prize *Prize) error

	Update(prize *Prize) error

	Delete(id string) error
}
