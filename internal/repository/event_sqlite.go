package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type EventSqlite struct {
	db *sql.DB
}

func NewEventSqlite(db *sql.DB) *EventSqlite {
	return &EventSqlite{db}
}

func (r *EventSqlite) FindAll() (events []internal.Event, err error) {
	rows, err := r.db.Query("SELECT e.`id`, e.`name`, e.`start_time_mili`, e.`end_time_mili`, e.`schedule_id` FROM `events` AS `e`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var event internal.Event
		err = rows.Scan(&event.ID, &event.Name, &event.Start_time_mili, &event.End_time_mili, &event.Schedule_id)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *EventSqlite) FindById(id string) (event internal.Event, err error) {
	row := r.db.QueryRow("SELECT e.`id`, e.`name`, e.`start_time_mili`, e.`end_time_mili`, e.`schedule_id` FROM `events` AS `e` WHERE e.`id` == ?;", id)

	err = row.Scan(&event.ID, &event.Name, &event.Start_time_mili, &event.End_time_mili, &event.Schedule_id)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrEventRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *EventSqlite) Save(event *internal.Event) (err error) {
	_, err = r.db.Exec("INSERT INTO `events` (`id`, `name`, `start_time_mili`, `end_time_mili`, `schedule_id`) VALUES (?, ?, ?, ?, ?)", event.ID, event.Name, event.Start_time_mili, event.End_time_mili)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrEventRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *EventSqlite) Update(event *internal.Event) (err error) {
	_, err = r.db.Exec("UPDATE `events` SET `name` = ?, `start_time_mili` = ?, `end_time_mili` = ?, `schedule_id` = ? WHERE `id` = ?;", event.Name, event.Start_time_mili, event.End_time_mili, event.Schedule_id, event.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrEventRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *EventSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `events` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrEventRepositoryNotFound
		return
	}

	return
}

func (r *EventSqlite) GetBasicInfo() (count internal.EventInfoCount, err error) {
	query := `
	SELECT 
    (SELECT COUNT(*) FROM schedules) AS schedules_count,
    (SELECT COUNT(*) FROM runs) AS runs_count,
    (SELECT COUNT(*) FROM prizes) AS prizes_count,
	(SELECT COUNT(*) FROM bids) AS bids_count,
	(SELECT COUNT(*) FROM donations) AS donations_count,
	(SELECT COUNT(*) FROM users) AS users_count;
	`
	row := r.db.QueryRow(query)
	err = row.Scan(&count.Schedules_count, &count.Runs_count, &count.Prizes_count, &count.Bids_count, &count.Donations_count, &count.Users_count)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrEventRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}
