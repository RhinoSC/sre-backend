package internal

import "errors"

type RunTeamPlayers struct {
	UserID string
	User   User
}

type RunTeams struct {
	ID      string
	Name    string
	Players []RunTeamPlayers
}

type RunMetadata struct {
	ID             string
	RunID          string
	Category       string
	Platform       string
	TwitchGameName string
	RunName        string
	Note           string
}

type Run struct {
	ID             string
	Name           string
	StartTimeMili  int64
	EstimateString string
	EstimateMili   int64
	SetupTimeMili  int64
	Status         string
	RunMetadata
	Teams      []RunTeams
	Bids       []Bid
	ScheduleId string
}

var (
	// REPOSITORY ERRORS
	ErrRunRepositoryNotFound   = errors.New("repository: run not found")
	ErrRunRepositoryDuplicated = errors.New("repository: run already exists")
	// SERVICE ERRORS
	ErrRunDatabase          = errors.New("database error")
	ErrRunServiceNotFound   = errors.New("service: user not found")
	ErrRunServiceDuplicated = errors.New("service: user already exists")
)

type RunRepository interface {
	FindAll() ([]Run, error)

	FindById(id string) (Run, error)

	Save(run *Run) error

	Update(run *Run) error

	Delete(id string) error

	UpdateRunOrder(runs []Run) error
}

type RunService interface {
	FindAll() ([]Run, error)

	FindById(id string) (Run, error)

	Save(run *Run) error

	Update(run *Run) error

	Delete(id string) error

	UpdateRunOrder(runs []Run) error
}
