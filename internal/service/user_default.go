package service

import (
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
