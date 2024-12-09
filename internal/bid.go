package internal

import (
	"database/sql"
	"errors"
)

type BidType string

const (
	Bidwar BidType = "bidwar"
	Total  BidType = "total"
	Goal   BidType = "goal"
)

type BidOptions struct {
	ID            string
	Name          string
	CurrentAmount float64
	BidID         string
}

type BidOptionsSQL struct {
	ID            sql.NullString
	Name          sql.NullString
	CurrentAmount sql.NullFloat64
	BidID         sql.NullString
}

type Bid struct {
	ID               string
	Bidname          string
	Goal             float64
	CurrentAmount    float64
	Description      string
	Type             BidType
	CreateNewOptions bool
	Status           string
	RunID            string
	BidOptions       []BidOptions
}

var (
	// REPOSITORY ERRORS
	ErrBidRepositoryNotFound   = errors.New("repository: bid not found")
	ErrBidRepositoryDuplicated = errors.New("repository: bid already exists")
	// SERVICE ERRORS
	ErrBidDatabase          = errors.New("database error")
	ErrBidServiceNotFound   = errors.New("service: bid not found")
	ErrBidServiceDuplicated = errors.New("service: bid already exists")
)

type BidRepository interface {
	FindAll() ([]Bid, error)

	FindById(id string) (Bid, error)

	Save(bid *Bid) error

	Update(bid *Bid) error

	Delete(id string) error
}

type BidService interface {
	FindAll() ([]Bid, error)

	FindById(id string) (Bid, error)

	Save(bid *Bid) error

	Update(bid *Bid) error

	Delete(id string) error
}
