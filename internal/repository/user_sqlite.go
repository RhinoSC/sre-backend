package repository

import (
	"database/sql"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db}
}

func (r *UserSqlite) FindAll() (users []internal.User, err error) {
	rows, err := r.db.Query("SELECT u.`id`, u.`name`, u.`username` FROM `users` AS `u`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var user internal.User
		err = rows.Scan(&user.ID, &user.Name, &user.Username)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}
