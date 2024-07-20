package internal

import "errors"

type Team struct {
	ID   string
	Name string
}

var (
	// REPOSITORY ERRORS
	ErrTeamRepositoryNotFound   = errors.New("repository: team not found")
	ErrTeamRepositoryDuplicated = errors.New("repository: team already exists")
	// SERVICE ERRORS
	ErrTeamDatabase          = errors.New("database error")
	ErrTeamServiceNotFound   = errors.New("service: team not found")
	ErrTeamServiceDuplicated = errors.New("service: team already exists")
)

type TeamRepository interface {
	FindAll() ([]Team, error)

	FindById(id string) (Team, error)

	Save(team *Team) error

	Update(team *Team) error

	Delete(id string) error
}

type TeamService interface {
	FindAll() ([]Team, error)

	FindById(id string) (Team, error)

	Save(team *Team) error

	Update(team *Team) error

	Delete(id string) error
}
