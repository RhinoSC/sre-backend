package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/RhinoSC/sre-backend/internal/utils"
	"github.com/mattn/go-sqlite3"
)

type AdminSqlite struct {
	db *sql.DB
}

func NewAdminSqlite(db *sql.DB) *AdminSqlite {
	return &AdminSqlite{db}
}

func (r *AdminSqlite) FindAll() (admins []internal.Admin, err error) {
	rows, err := r.db.Query("SELECT a.`id`, a.`username`, a.`password` FROM `admins` AS a;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var admin internal.Admin
		err = rows.Scan(&admin.ID, &admin.Username, &admin.Password)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		admins = append(admins, admin)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *AdminSqlite) FindById(id string) (admin internal.Admin, err error) {
	row := r.db.QueryRow("SELECT  a.`id`, a.`username`, a.`password` FROM `admins` AS `a` WHERE a.`id` == ?;", id)

	err = row.Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrAdminRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *AdminSqlite) Save(admin *internal.Admin) (err error) {
	_, err = r.db.Exec("INSERT INTO `admins` (`id`, `username`, password) VALUES (?, ?, ?)", admin.ID, admin.Username, admin.Password)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrAdminRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *AdminSqlite) Update(admin *internal.Admin) (err error) {
	_, err = r.db.Exec("UPDATE `admins` SET `username` = ?, `password` = ? WHERE `id` = ?;", admin.Username, admin.Password, admin.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrAdminRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *AdminSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `admins` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrAdminRepositoryNotFound
		return
	}

	return
}

func (r *AdminSqlite) Login(username string, password string) (admin internal.Admin, err error) {
	row := r.db.QueryRow("SELECT  a.`id`, a.`username`, a.`password` FROM `admins` AS `a` WHERE a.`username` == ?;", username)

	err = row.Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrAdminRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	checkPassword := utils.CheckPasswordHash(password, admin.Password)
	if !checkPassword {
		err = internal.ErrAdminRepositoryInvalidPassword
		logger.Log.Error(err.Error())
		return
	}

	return
}
