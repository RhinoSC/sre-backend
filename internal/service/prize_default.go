package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type PrizeDefault struct {
	rp internal.PrizeRepository
}

func NewPrizeDefault(rp internal.PrizeRepository) *PrizeDefault {
	return &PrizeDefault{
		rp: rp,
	}
}

func (s *PrizeDefault) FindAll() (prizes []internal.Prize, err error) {
	prizes, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all prizes: %w", err)
		return
	}
	return
}

func (s *PrizeDefault) FindById(id string) (prize internal.Prize, err error) {
	prize, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrPrizeRepositoryNotFound):
			err = fmt.Errorf("error finding prize by id: %w", internal.ErrPrizeServiceNotFound)
		default:
			err = fmt.Errorf("error finding prize by id: %w", err)
		}
		return
	}
	return
}

func (s *PrizeDefault) Save(prize *internal.Prize) (err error) {
	err = s.rp.Save(prize)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrPrizeRepositoryDuplicated):
			err = fmt.Errorf("error saving prize: %w", internal.ErrPrizeServiceDuplicated)
		default:
			err = fmt.Errorf("error saving prize: %w", err)
		}
		return
	}

	return
}

func (s *PrizeDefault) Update(prize *internal.Prize) (err error) {
	err = s.rp.Update(prize)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrPrizeRepositoryDuplicated):
			err = fmt.Errorf("error updating prize: %w", internal.ErrPrizeServiceDuplicated)
		default:
			err = fmt.Errorf("error updating prize: %w", err)
		}
		return
	}
	return
}

func (s *PrizeDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrPrizeRepositoryNotFound):
			err = fmt.Errorf("error deleting prize: %w", internal.ErrPrizeServiceNotFound)
		default:
			err = fmt.Errorf("error deleting prize: %w", err)
		}
		return
	}
	return
}
