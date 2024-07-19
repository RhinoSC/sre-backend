package internal

import "errors"

type Event struct {
	ID              string
	Name            string
	Start_time_mili int64
	End_time_mili   int64
}

var (
	// REPOSITORY ERRORS
	ErrEventRepositoryNotFound   = errors.New("repository: event not found")
	ErrEventRepositoryDuplicated = errors.New("repository: event already exists")
	// SERVICE ERRORS
	ErrEventDatabase          = errors.New("database error")
	ErrEventServiceNotFound   = errors.New("service: event not found")
	ErrEventServiceDuplicated = errors.New("service: event already exists")
)

type EventRepository interface {
	FindAll() ([]Event, error)

	FindById(id string) (Event, error)

	Save(event *Event) error

	Update(event *Event) error

	Delete(id string) error
}

type EventService interface {
	FindAll() ([]Event, error)

	FindById(id string) (Event, error)

	Save(event *Event) error

	Update(event *Event) error

	Delete(id string) error
}
