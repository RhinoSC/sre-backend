package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type TeamDefault struct {
	rp internal.TeamRepository
}

func NewTeamDefault(rp internal.TeamRepository) *TeamDefault {
	return &TeamDefault{
		rp: rp,
	}
}

func (s *TeamDefault) FindAll() (teams []internal.Team, err error) {
	teams, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all teams: %w", err)
		return
	}
	return
}

func (s *TeamDefault) FindById(id string) (team internal.Team, err error) {
	team, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrTeamRepositoryNotFound):
			err = fmt.Errorf("error finding team by id: %w", internal.ErrTeamServiceNotFound)
		default:
			err = fmt.Errorf("error finding team by id: %w", err)
		}
		return
	}
	return
}

func (s *TeamDefault) Save(team *internal.Team) (err error) {
	err = s.rp.Save(team)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrTeamRepositoryDuplicated):
			err = fmt.Errorf("error saving team: %w", internal.ErrTeamServiceDuplicated)
		default:
			err = fmt.Errorf("error saving team: %w", err)
		}
		return
	}

	return
}

func (s *TeamDefault) Update(team *internal.Team) (err error) {
	err = s.rp.Update(team)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrTeamRepositoryDuplicated):
			err = fmt.Errorf("error updating team: %w", internal.ErrTeamServiceDuplicated)
		default:
			err = fmt.Errorf("error updating team: %w", err)
		}
		return
	}
	return
}

func (s *TeamDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrTeamRepositoryNotFound):
			err = fmt.Errorf("error deleting team: %w", internal.ErrTeamServiceNotFound)
		default:
			err = fmt.Errorf("error deleting team: %w", err)
		}
		return
	}
	return
}
