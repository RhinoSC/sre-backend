package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type DonationDefault struct {
	rp internal.DonationRepository
}

func NewDonationDefault(rp internal.DonationRepository) *DonationDefault {
	return &DonationDefault{
		rp: rp,
	}
}

func (s *DonationDefault) FindAll() (donations []internal.Donation, err error) {
	donations, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all donations: %w", err)
		return
	}
	return
}

func (s *DonationDefault) FindAllWithBidDetails() (donations []internal.DonationWithBidDetails, err error) {
	donations, err = s.rp.FindAllWithBidDetails()
	if err != nil {
		err = fmt.Errorf("error finding all donations: %w", err)
		return
	}
	return
}

func (s *DonationDefault) FindById(id string) (donation internal.Donation, err error) {
	donation, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryNotFound):
			err = fmt.Errorf("error finding donation by id: %w", internal.ErrDonationServiceNotFound)
		default:
			err = fmt.Errorf("error finding donation by id: %w", err)
		}
		return
	}
	return
}

func (s *DonationDefault) FindByIdWithBidDetails(id string) (donation internal.DonationWithBidDetails, err error) {
	donation, err = s.rp.FindByIdWithBidDetails(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryNotFound):
			err = fmt.Errorf("error finding donation with bid details by id: %w", internal.ErrDonationServiceNotFound)
		default:
			err = fmt.Errorf("error finding donation with bid details by id: %w", err)
		}
		return
	}
	return
}

func (s *DonationDefault) FindByEventID(id string) (donations []internal.Donation, err error) {
	donations, err = s.rp.FindByEventID(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryNotFound):
			err = fmt.Errorf("error finding donation by event id: %w", internal.ErrDonationServiceNotFound)
		default:
			err = fmt.Errorf("error finding donation by event id: %w", err)
		}
		return
	}
	return
}

func (s *DonationDefault) FindByEventIDWithBidDetails(id string) (donations []internal.DonationWithBidDetails, err error) {
	donations, err = s.rp.FindByEventIDWithBidDetails(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryNotFound):
			err = fmt.Errorf("error finding donation by event id: %w", internal.ErrDonationServiceNotFound)
		default:
			err = fmt.Errorf("error finding donation by event id: %w", err)
		}
		return
	}
	return
}

func (s *DonationDefault) Save(donation *internal.DonationWithBidDetails) (err error) {
	err = s.rp.Save(donation)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryDuplicated):
			err = fmt.Errorf("error saving donation: %w", internal.ErrDonationServiceDuplicated)
		default:
			err = fmt.Errorf("error saving donation: %w", err)
		}
		return
	}

	return
}

func (s *DonationDefault) Update(donation *internal.DonationWithBidDetails) (err error) {
	err = s.rp.Update(donation)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryDuplicated):
			err = fmt.Errorf("error updating donation: %w", internal.ErrDonationServiceDuplicated)
		default:
			err = fmt.Errorf("error updating donation: %w", err)
		}
		return
	}
	return
}

func (s *DonationDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrDonationRepositoryNotFound):
			err = fmt.Errorf("error deleting donation: %w", internal.ErrDonationServiceNotFound)
		default:
			err = fmt.Errorf("error deleting donation: %w", err)
		}
		return
	}
	return
}
