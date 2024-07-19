package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type EventDefault struct {
	rp internal.EventRepository
}

func NewEventDefault(rp internal.EventRepository) *EventDefault {
	return &EventDefault{
		rp: rp,
	}
}

func (s *EventDefault) FindAll() (events []internal.Event, err error) {
	events, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all events: %w", err)
		return
	}
	return
}

func (s *EventDefault) FindById(id string) (event internal.Event, err error) {
	event, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEventRepositoryNotFound):
			err = fmt.Errorf("error finding event by id: %w", internal.ErrEventServiceNotFound)
		default:
			err = fmt.Errorf("error finding event by id: %w", err)
		}
		return
	}
	return
}

func (s *EventDefault) Save(event *internal.Event) (err error) {
	err = s.rp.Save(event)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEventRepositoryDuplicated):
			err = fmt.Errorf("error saving event: %w", internal.ErrEventServiceDuplicated)
		default:
			err = fmt.Errorf("error saving event: %w", err)
		}
		return
	}

	return
}

func (s *EventDefault) Update(event *internal.Event) (err error) {
	err = s.rp.Update(event)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEventRepositoryDuplicated):
			err = fmt.Errorf("error updating event: %w", internal.ErrEventServiceDuplicated)
		default:
			err = fmt.Errorf("error updating event: %w", err)
		}
		return
	}
	return
}

func (s *EventDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEventRepositoryNotFound):
			err = fmt.Errorf("error deleting event: %w", internal.ErrEventServiceNotFound)
		default:
			err = fmt.Errorf("error deleting event: %w", err)
		}
		return
	}
	return
}
