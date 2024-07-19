package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type UserDefault struct {
	rp internal.UserRepository
}

func NewUserDefault(rp internal.UserRepository) *UserDefault {
	return &UserDefault{
		rp: rp,
	}
}

func (s *UserDefault) FindAll() (users []internal.User, err error) {
	users, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all users: %w", err)
		return
	}
	return
}

func (s *UserDefault) FindById(id string) (user internal.User, err error) {
	user, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrUserRepositoryNotFound):
			err = fmt.Errorf("error finding user by id: %w", internal.ErrUserServiceNotFound)
		default:
			err = fmt.Errorf("error finding user by id: %w", err)
		}
		return
	}
	return
}

func (s *UserDefault) FindByUsername(username string) (user internal.User, err error) {
	user, err = s.rp.FindByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrUserRepositoryNotFound):
			err = fmt.Errorf("error finding user by username: %w", internal.ErrUserServiceNotFound)
		default:
			err = fmt.Errorf("error finding user by username: %w", err)
		}
		return
	}
	return
}

func (s *UserDefault) Save(user *internal.User) (err error) {
	err = s.rp.Save(user)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrUserRepositoryDuplicated):
			err = fmt.Errorf("error saving user: %w", internal.ErrUserServiceDuplicated)
		default:
			err = fmt.Errorf("error saving user: %w", err)
		}
		return
	}

	return
}

func (s *UserDefault) Update(user *internal.User) (err error) {
	err = s.rp.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrUserRepositoryDuplicated):
			err = fmt.Errorf("error updating user: %w", internal.ErrUserServiceDuplicated)
		default:
			err = fmt.Errorf("error updating user: %w", err)
		}
		return
	}
	return
}
