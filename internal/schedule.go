package internal

import "errors"

type Schedule struct {
	ID              string
	Name            string
	Start_time_mili int64
	End_time_mili   int64
	EventID         string
}

var (
	// REPOSITORY ERRORS
	ErrScheduleRepositoryNotFound   = errors.New("repository: schedule not found")
	ErrScheduleRepositoryDuplicated = errors.New("repository: schedule already exists")
	// SERVICE ERRORS
	ErrScheduleDatabase          = errors.New("database error")
	ErrScheduleServiceNotFound   = errors.New("service: schedule not found")
	ErrScheduleServiceDuplicated = errors.New("service: schedule already exists")
)

type ScheduleRepository interface {
	FindAll() ([]Schedule, error)

	FindById(id string) (Schedule, error)

	Save(schedule *Schedule) error

	Update(schedule *Schedule) error

	Delete(id string) error
}

type ScheduleService interface {
	FindAll() ([]Schedule, error)

	FindById(id string) (Schedule, error)

	Save(schedule *Schedule) error

	Update(schedule *Schedule) error

	Delete(id string) error
}
