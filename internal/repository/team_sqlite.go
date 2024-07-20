package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type TeamSqlite struct {
	db *sql.DB
}

func NewTeamSqlite(db *sql.DB) *TeamSqlite {
	return &TeamSqlite{db}
}

func (r *TeamSqlite) FindAll() (teams []internal.Team, err error) {
	rows, err := r.db.Query("SELECT t.`id`, t.`name` FROM `teams` AS t;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var team internal.Team
		err = rows.Scan(&team.ID, &team.Name)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		teams = append(teams, team)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *TeamSqlite) FindById(id string) (team internal.Team, err error) {
	row := r.db.QueryRow("SELECT t.`id`, t.`name` FROM `teams` AS `t` WHERE t.`id` == ?;", id)

	err = row.Scan(&team.ID, &team.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrTeamRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *TeamSqlite) Save(team *internal.Team) (err error) {
	_, err = r.db.Exec("INSERT INTO `teams` (`id`, `name`) VALUES (?, ?)", team.ID, team.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrTeamRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *TeamSqlite) Update(team *internal.Team) (err error) {
	_, err = r.db.Exec("UPDATE `teams` SET `name` = ? WHERE `id` = ?;", team.ID, team.Name)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrTeamRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *TeamSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `teams` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrTeamRepositoryNotFound
		return
	}

	return
}
