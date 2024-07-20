package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type BidDefault struct {
	rp internal.BidRepository
}

func NewBidDefault(rp internal.BidRepository) *BidDefault {
	return &BidDefault{
		rp: rp,
	}
}

func (s *BidDefault) FindAll() (bids []internal.Bid, err error) {
	bids, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all bids: %w", err)
		return
	}
	return
}

func (s *BidDefault) FindById(id string) (bid internal.Bid, err error) {
	bid, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrBidRepositoryNotFound):
			err = fmt.Errorf("error finding bid by id: %w", internal.ErrBidServiceNotFound)
		default:
			err = fmt.Errorf("error finding bid by id: %w", err)
		}
		return
	}
	return
}

func (s *BidDefault) Save(bid *internal.Bid) (err error) {
	err = s.rp.Save(bid)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrBidRepositoryDuplicated):
			err = fmt.Errorf("error saving bid: %w", internal.ErrBidServiceDuplicated)
		default:
			err = fmt.Errorf("error saving bid: %w", err)
		}
		return
	}

	return
}

func (s *BidDefault) Update(bid *internal.Bid) (err error) {
	err = s.rp.Update(bid)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrBidRepositoryDuplicated):
			err = fmt.Errorf("error updating bid: %w", internal.ErrBidServiceDuplicated)
		default:
			err = fmt.Errorf("error updating bid: %w", err)
		}
		return
	}
	return
}

func (s *BidDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrBidRepositoryNotFound):
			err = fmt.Errorf("error deleting bid: %w", internal.ErrBidServiceNotFound)
		default:
			err = fmt.Errorf("error deleting bid: %w", err)
		}
		return
	}
	return
}
