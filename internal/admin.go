package internal

import "errors"

type Admin struct {
	User
	Password string
}

var (
	// REPOSITORY ERRORS
	ErrAdminRepositoryNotFound        = errors.New("repository: admin not found")
	ErrAdminRepositoryDuplicated      = errors.New("repository: admin already exists")
	ErrAdminRepositoryInvalidPassword = errors.New("repository: invalid password")
	// SERVICE ERRORS
	ErrAdminDatabase               = errors.New("database error")
	ErrAdminServiceNotFound        = errors.New("service: admin not found")
	ErrAdminServiceDuplicated      = errors.New("service: admin already exists")
	ErrAdminServiceInvalidPassword = errors.New("service: invalid password")
)

type AdminRepository interface {
	FindAll() ([]Admin, error)

	FindById(id string) (Admin, error)

	Save(admin *Admin) error

	Update(admin *Admin) error

	Delete(id string) error

	Login(username string, password string) (Admin, error)
}

type AdminService interface {
	FindAll() ([]Admin, error)

	FindById(id string) (Admin, error)

	Save(admin *Admin) error

	Update(admin *Admin) error

	Delete(id string) error

	Login(username string, password string) (Admin, error)
}
