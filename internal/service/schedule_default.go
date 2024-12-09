package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type ScheduleDefault struct {
	rp internal.ScheduleRepository
}

func NewScheduleDefault(rp internal.ScheduleRepository) *ScheduleDefault {
	return &ScheduleDefault{
		rp: rp,
	}
}

func (s *ScheduleDefault) FindAll() (schedules []internal.Schedule, err error) {
	schedules, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all schedules: %w", err)
		return
	}
	return
}

func (s *ScheduleDefault) FindById(id string) (schedule internal.Schedule, err error) {
	schedule, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrScheduleRepositoryNotFound):
			err = fmt.Errorf("error finding schedule by id: %w", internal.ErrScheduleServiceNotFound)
		default:
			err = fmt.Errorf("error finding schedule by id: %w", err)
		}
		return
	}
	return
}

func (s *ScheduleDefault) Save(schedule *internal.Schedule) (err error) {
	err = s.rp.Save(schedule)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrScheduleRepositoryDuplicated):
			err = fmt.Errorf("error saving schedule: %w", internal.ErrScheduleServiceDuplicated)
		default:
			err = fmt.Errorf("error saving schedule: %w", err)
		}
		return
	}

	return
}

func (s *ScheduleDefault) Update(schedule *internal.Schedule) (err error) {
	err = s.rp.Update(schedule)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrScheduleRepositoryDuplicated):
			err = fmt.Errorf("error updating schedule: %w", internal.ErrScheduleServiceDuplicated)
		default:
			err = fmt.Errorf("error updating schedule: %w", err)
		}
		return
	}
	return
}

func (s *ScheduleDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrScheduleRepositoryNotFound):
			err = fmt.Errorf("error deleting schedule: %w", internal.ErrScheduleServiceNotFound)
		default:
			err = fmt.Errorf("error deleting schedule: %w", err)
		}
		return
	}
	return
}
