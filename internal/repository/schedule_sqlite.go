package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type ScheduleSqlite struct {
	db *sql.DB
}

func NewScheduleSqlite(db *sql.DB) *ScheduleSqlite {
	return &ScheduleSqlite{db}
}

func (r *ScheduleSqlite) FindAll() (schedules []internal.Schedule, err error) {
	rows, err := r.db.Query("SELECT s.`id`, s.`name`, s.`start_time_mili`, s.`end_time_mili`, s.`event_id` FROM `schedules` AS `s`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var schedule internal.Schedule
		err = rows.Scan(&schedule.ID, &schedule.Name, &schedule.Start_time_mili, &schedule.End_time_mili, &schedule.EventID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		schedules = append(schedules, schedule)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *ScheduleSqlite) FindById(id string) (schedule internal.Schedule, err error) {
	row := r.db.QueryRow("SELECT s.`id`, s.`name`, s.`start_time_mili`, s.`end_time_mili`, s.`event_id` FROM `schedules` AS `s` WHERE s.`id` == ?;", id)

	err = row.Scan(&schedule.ID, &schedule.Name, &schedule.Start_time_mili, &schedule.End_time_mili, &schedule.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrScheduleRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *ScheduleSqlite) Save(schedule *internal.Schedule) (err error) {
	_, err = r.db.Exec("INSERT INTO `schedules` (`id`, `name`, `start_time_mili`, `end_time_mili`, `event_id`) VALUES (?, ?, ?, ?, ?)", schedule.ID, schedule.Name, schedule.Start_time_mili, schedule.End_time_mili, schedule.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrScheduleRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *ScheduleSqlite) Update(schedule *internal.Schedule) (err error) {
	_, err = r.db.Exec("UPDATE `schedules` SET `name` = ?, `start_time_mili` = ?, `end_time_mili` = ?, `event_id` = ? WHERE `id` = ?;", schedule.Name, schedule.Start_time_mili, schedule.End_time_mili, schedule.EventID, schedule.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrScheduleRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *ScheduleSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `schedules` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrScheduleRepositoryNotFound
		return
	}

	return
}
