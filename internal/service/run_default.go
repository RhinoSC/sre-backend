package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type RunDefault struct {
	rp internal.RunRepository
}

func NewRunDefault(rp internal.RunRepository) *RunDefault {
	return &RunDefault{rp}
}

func (s *RunDefault) FindAll() (runs []internal.Run, err error) {
	runs, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all runs: %w", err)
		return
	}

	return
}

func (s *RunDefault) FindById(id string) (run internal.Run, err error) {
	run, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRunRepositoryNotFound):
			err = fmt.Errorf("error finding run by id: %w", internal.ErrRunServiceNotFound)
		default:
			err = fmt.Errorf("error finding run by id: %w", err)
		}
		return
	}
	return
}

func (s *RunDefault) Save(run *internal.Run) (err error) {
	err = s.rp.Save(run)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRunRepositoryDuplicated):
			err = fmt.Errorf("error saving run: %w", internal.ErrRunServiceDuplicated)
		default:
			err = fmt.Errorf("error saving run: %w", err)
		}
		return
	}
	return
}

func (s *RunDefault) Update(run *internal.Run) (err error) {
	err = s.rp.Update(run)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRunRepositoryDuplicated):
			err = fmt.Errorf("error updating run: %w", internal.ErrRunServiceDuplicated)
		default:
			err = fmt.Errorf("error updating run: %w", err)
		}
		return
	}
	return
}

func (s *RunDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrRunRepositoryNotFound):
			err = fmt.Errorf("error deleting run: %w", internal.ErrRunServiceNotFound)
		default:
			err = fmt.Errorf("error deleting run: %w", err)
		}
		return
	}
	return
}

func (s *RunDefault) UpdateRunOrder(runs []internal.Run) (err error) {
	err = s.rp.UpdateRunOrder(runs)
	if err != nil {
		switch {
		default:
			err = fmt.Errorf("error updating runs order: %w", err)
		}
		return
	}
	return
}
