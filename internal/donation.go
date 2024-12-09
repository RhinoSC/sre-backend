package internal

import (
	"database/sql"
	"errors"
)

type Donation struct {
	ID          string
	Name        string
	Email       string
	TimeMili    int64
	Amount      float64
	Description string
	ToBid       bool
	EventID     string
}

type DonationWithBidDetails struct {
	Donation
	BidDetails    *DonationBidDetails
	NewBidDetails *DonationBidDetails
}

type DonationBidDetails struct {
	BidID            string
	Bidname          string
	Goal             float64
	CurrentAmount    float64
	BidDescription   string
	Type             BidType
	CreateNewOptions bool
	RunID            string
	OptionID         string
	OptionName       string
	OptionAmount     float64
}

type DonationBidDetailsDB struct {
	Donation
	BidID            sql.NullString
	Bidname          sql.NullString
	Goal             sql.NullFloat64
	CurrentAmount    sql.NullFloat64
	BidDescription   sql.NullString
	Type             sql.NullString
	CreateNewOptions sql.NullBool
	RunID            sql.NullString
	OptionID         sql.NullString
	OptionName       sql.NullString
	OptionAmount     sql.NullFloat64
}

type DonationWithBidDetailsDB struct {
	Donation
	BidDetails    *DonationBidDetailsDB
	NewBidDetails *DonationBidDetailsDB
}

var (
	// REPOSITORY ERRORS
	ErrDonationRepositoryNotFound   = errors.New("repository: donation not found")
	ErrDonationRepositoryDuplicated = errors.New("repository: donation already exists")
	// SERVICE ERRORS
	ErrDonationDatabase          = errors.New("database error")
	ErrDonationServiceNotFound   = errors.New("service: donation not found")
	ErrDonationServiceDuplicated = errors.New("service: donation already exists")
)

type DonationRepository interface {
	FindAll() ([]Donation, error)

	FindAllWithBidDetails() ([]DonationWithBidDetails, error)

	FindById(id string) (Donation, error)

	FindByIdWithBidDetails(id string) (DonationWithBidDetails, error)

	FindByEventID(id string) ([]Donation, error)

	FindByEventIDWithBidDetails(id string) ([]DonationWithBidDetails, error)

	FindTotalDonatedByEventID(id string) (float64, error)

	Save(donation *DonationWithBidDetails) error

	Update(donation *DonationWithBidDetails) error

	Delete(id string) error
}

type DonationService interface {
	FindAll() ([]Donation, error)

	FindAllWithBidDetails() ([]DonationWithBidDetails, error)

	FindById(id string) (Donation, error)

	FindByIdWithBidDetails(id string) (DonationWithBidDetails, error)

	FindByEventID(id string) ([]Donation, error)

	FindByEventIDWithBidDetails(id string) ([]DonationWithBidDetails, error)

	FindTotalDonatedByEventID(id string) (float64, error)

	Save(donation *DonationWithBidDetails) error

	Update(donation *DonationWithBidDetails) error

	Delete(id string) error
}
