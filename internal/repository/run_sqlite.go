package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type RunSqlite struct {
	db *sql.DB
}

func NewRunSqlite(db *sql.DB) *RunSqlite {
	return &RunSqlite{db}
}

func (r *RunSqlite) FindAll() (runs []internal.Run, err error) {
	rows, err := r.db.Query("SELECT `r`.id, `r`.name, `r`.start_time_mili, `r`.estimate_string, `r`.estimate_mili, `rm`.category, `rm`.platform, `rm`.twitch_game_name, `rm`.run_name, `rm`.note, `r`.schedule_id FROM runs AS `r` JOIN run_metadata AS `rm` ON `r`.id = `rm`.run_id;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var run internal.Run
		err = rows.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.RunName, &run.RunMetadata.Note, &run.ScheduleId)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		runs = append(runs, run)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) FindById(id string) (run internal.Run, err error) {
	row := r.db.QueryRow("SELECT `r`.id, `r`.name, `r`.start_time_mili, `r`.estimate_string, `r`.estimate_mili, `rm`.category, `rm`.platform, `rm`.twitch_game_name, `rm`.run_name, `rm`.note, `r`.schedule_id FROM runs AS `r` JOIN run_metadata AS `rm` ON `r`.id = `rm`.run_id WHERE `r`.id = ?;", id)
	err = row.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.RunName, &run.RunMetadata.Note, &run.ScheduleId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrRunRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) Save(run *internal.Run) (err error) {
	_, err = r.db.Exec("INSERT INTO runs (id, name, start_time_mili, estimate_string, estimate_mili, schedule_id) VALUES (?, ?, ?, ?, ?, ?);", run.ID, run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.ScheduleId)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch {
			case sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique:
				err = internal.ErrRunRepositoryDuplicated
			default:
				err = internal.ErrRunDatabase
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	_, err = r.db.Exec("INSERT INTO run_metadata (id, run_id, category, platform, twitch_game_name, run_name, note) VALUES (?, ?, ?, ?, ?, ?, ?);", run.ID, run.ID, run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.RunName, run.RunMetadata.Note)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch {
			case sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique:
				err = internal.ErrRunRepositoryDuplicated
			default:
				err = internal.ErrRunDatabase
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *RunSqlite) Update(run *internal.Run) (err error) {
	_, err = r.db.Exec("UPDATE runs SET name = ?, start_time_mili = ?, estimate_string = ?, estimate_mili = ?, schedule_id = ? WHERE id = ?;", run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.ScheduleId, run.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch {
			case sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique:
				err = internal.ErrRunRepositoryDuplicated
			default:
				err = internal.ErrRunDatabase
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	_, err = r.db.Exec("UPDATE run_metadata SET category = ?, platform = ?, twitch_game_name = ?, run_name = ?, note = ? WHERE run_id = ?;", run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.RunName, run.RunMetadata.Note, run.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch {
			case sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique:
				err = internal.ErrRunRepositoryDuplicated
			default:
				err = internal.ErrRunDatabase
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *RunSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM runs WHERE id = ?;", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	if rowsAffected == 0 {
		err = internal.ErrRunRepositoryNotFound
		return
	}

	return
}
