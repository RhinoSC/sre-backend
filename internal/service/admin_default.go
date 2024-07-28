package service

import (
	"errors"
	"fmt"

	"github.com/RhinoSC/sre-backend/internal"
)

type AdminDefault struct {
	rp internal.AdminRepository
}

func NewAdminDefault(rp internal.AdminRepository) *AdminDefault {
	return &AdminDefault{
		rp: rp,
	}
}

func (s *AdminDefault) FindAll() (admins []internal.Admin, err error) {
	admins, err = s.rp.FindAll()
	if err != nil {
		err = fmt.Errorf("error finding all admins: %w", err)
		return
	}
	return
}

func (s *AdminDefault) FindById(id string) (admin internal.Admin, err error) {
	admin, err = s.rp.FindById(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrAdminRepositoryNotFound):
			err = fmt.Errorf("error finding admin by id: %w", internal.ErrAdminServiceNotFound)
		default:
			err = fmt.Errorf("error finding admin by id: %w", err)
		}
		return
	}
	return
}

func (s *AdminDefault) Save(admin *internal.Admin) (err error) {
	err = s.rp.Save(admin)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrAdminRepositoryDuplicated):
			err = fmt.Errorf("error saving admin: %w", internal.ErrAdminServiceDuplicated)
		default:
			err = fmt.Errorf("error saving admin: %w", err)
		}
		return
	}

	return
}

func (s *AdminDefault) Update(admin *internal.Admin) (err error) {
	err = s.rp.Update(admin)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrAdminRepositoryDuplicated):
			err = fmt.Errorf("error updating admin: %w", internal.ErrAdminServiceDuplicated)
		default:
			err = fmt.Errorf("error updating admin: %w", err)
		}
		return
	}
	return
}

func (s *AdminDefault) Delete(id string) (err error) {
	err = s.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrAdminRepositoryNotFound):
			err = fmt.Errorf("error deleting admin: %w", internal.ErrAdminServiceNotFound)
		default:
			err = fmt.Errorf("error deleting admin: %w", err)
		}
		return
	}
	return
}

func (s *AdminDefault) Login(username string, password string) (admin internal.Admin, err error) {
	admin, err = s.rp.Login(username, password)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrAdminRepositoryNotFound):
			err = fmt.Errorf("error login. Admin not found: %w", internal.ErrAdminServiceNotFound)
		case errors.Is(err, internal.ErrAdminRepositoryInvalidPassword):
			err = fmt.Errorf("error login admin with invalid password: %w", internal.ErrAdminServiceInvalidPassword)
		default:
			err = fmt.Errorf("error login admin: %w", err)
		}
		return
	}
	return
}
