package internal

import "errors"

type UserSocials struct {
	Twitch   string
	Twitter  string
	Youtube  string
	Facebook string
}

type User struct {
	ID       string
	Name     string
	Username string
	UserSocials
}

var (
	// REPOSITORY ERRORS
	ErrUserRepositoryNotFound   = errors.New("repository: user not found")
	ErrUserRepositoryDuplicated = errors.New("repository: user already exists")
	// SERVICE ERRORS
	ErrUserDatabase          = errors.New("database error")
	ErrUserServiceNotFound   = errors.New("service: user not found")
	ErrUserServiceDuplicated = errors.New("service: user already exists")
)

type UserRepository interface {
	FindAll() ([]User, error)

	FindById(id string) (User, error)

	FindByUsername(username string) (User, error)

	Save(user *User) error

	Update(user *User) error

	Delete(id string) error
}

type UserService interface {
	FindAll() ([]User, error)

	FindById(id string) (User, error)

	FindByUsername(username string) (User, error)

	Save(user *User) error

	Update(user *User) error

	Delete(id string) error
}
