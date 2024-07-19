package internal

import "errors"

type User struct {
	ID       int
	Name     string
	Username string
}

var (
	// REPOSITORY ERRORS
	ErrUserRepositoryNotFound   = errors.New("repository: user not found")
	ErrUserRepositoryDuplicated = errors.New("repository: user already exists")
	// SERVICE ERRORS
	ErrUserDatabase          = errors.New("database error")
	ErrUserServiceDuplicated = errors.New("service: user already exists")
)

type UserRepository interface {
	FindAll() ([]User, error)

	FindById() (User, error)

	FindByUsername() (User, error)

	Save(user *User) error

	Update(user *User) error

	Delete(id string) error
}

type UserService interface {
	FindAll() ([]User, error)

	FindById() (User, error)

	FindByUsername() (User, error)

	Save(user *User) error

	Update(user *User) error

	Delete(id string) error
}
