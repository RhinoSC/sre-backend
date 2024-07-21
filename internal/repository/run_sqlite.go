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

// func (r *RunSqlite) FindAll() (runs []internal.Run, err error) {
// 	rows, err := r.db.Query("SELECT `r`.id, `r`.name, `r`.start_time_mili, `r`.estimate_string, `r`.estimate_mili, `rm`.category, `rm`.platform, `rm`.twitch_game_name, `rm`.run_name, `rm`.note, `r`.schedule_id FROM runs AS `r` JOIN run_metadata AS `rm` ON `r`.id = `rm`.run_id;")
// 	if err != nil {
// 		logger.Log.Error(err.Error())
// 		return
// 	}

// 	for rows.Next() {
// 		var run internal.Run
// 		err = rows.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.RunName, &run.RunMetadata.Note, &run.ScheduleId)
// 		if err != nil {
// 			logger.Log.Error(err.Error())
// 			return
// 		}

// 		runs = append(runs, run)
// 	}

// 	err = rows.Err()
// 	if err != nil {
// 		logger.Log.Error(err.Error())
// 		return
// 	}

// 	return
// }

// func (r *RunSqlite) FindById(id string) (run internal.Run, err error) {
// 	row := r.db.QueryRow("SELECT `r`.id, `r`.name, `r`.start_time_mili, `r`.estimate_string, `r`.estimate_mili, `rm`.category, `rm`.platform, `rm`.twitch_game_name, `rm`.run_name, `rm`.note, `r`.schedule_id FROM runs AS `r` JOIN run_metadata AS `rm` ON `r`.id = `rm`.run_id WHERE `r`.id = ?;", id)
// 	err = row.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.RunName, &run.RunMetadata.Note, &run.ScheduleId)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			err = internal.ErrRunRepositoryNotFound
// 		}
// 		logger.Log.Error(err.Error())
// 		return
// 	}

// 	return
// }

func (r *RunSqlite) FindAll() (runs []internal.Run, err error) {
	query := `
	SELECT 
		r.id AS run_id, r.name AS run_name, r.start_time_mili, r.estimate_string, r.estimate_mili, rm.category, rm.platform, rm.twitch_game_name, rm.note, r.schedule_id, 
		t.id AS team_id, t.name AS team_name, 
		u.id AS user_id, u.name AS user_name, u.username AS user_username, 
		um.id AS user_socials_id, um.twitch AS user_twitch, um.twitter AS user_twitter, um.youtube AS user_youtube, um.facebook AS user_facebook
	FROM 
		runs AS r 
		JOIN run_metadata AS rm ON r.id = rm.run_id 
		LEFT JOIN teams_runs AS tr ON r.id = tr.run_id 
		LEFT JOIN teams AS t ON tr.team_id = t.id 
		LEFT JOIN players AS pl ON t.id = pl.team_id 
		LEFT JOIN users AS u ON pl.user_id = u.id 
		LEFT JOIN user_socials AS um ON u.id = um.user_id;
	`
	rows, err := r.db.Query(query)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer rows.Close()

	runMap := make(map[string]*internal.Run)
	for rows.Next() {
		var run internal.Run
		var team internal.RunTeams
		var player internal.RunTeamPlayers
		var runID, teamID, userID sql.NullString
		var runName, teamName, userName, userUsername, socialsID, twitch, twitter, youtube, facebook sql.NullString
		var startTimeMili, estimateMili sql.NullInt64
		var estimateString, category, platform, twitchGameName, note, scheduleID sql.NullString

		err = rows.Scan(
			&runID, &runName, &startTimeMili, &estimateString, &estimateMili, &category, &platform, &twitchGameName, &note, &scheduleID,
			&teamID, &teamName,
			&userID, &userName, &userUsername,
			&socialsID, &twitch, &twitter, &youtube, &facebook,
		)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		runPtr, runExists := runMap[runID.String]
		if !runExists {
			run = internal.Run{
				ID:             runID.String,
				Name:           runName.String,
				StartTimeMili:  startTimeMili.Int64,
				EstimateString: estimateString.String,
				EstimateMili:   estimateMili.Int64,
				RunMetadata: internal.RunMetadata{
					ID:             runID.String,
					RunID:          runID.String,
					Category:       category.String,
					Platform:       platform.String,
					TwitchGameName: twitchGameName.String,
					Note:           note.String,
				},
				ScheduleId: scheduleID.String,
			}
			runMap[runID.String] = &run
		} else {
			run = *runPtr
		}

		if teamID.Valid {
			team.ID = teamID.String
			team.Name = teamName.String

			if userID.Valid {
				player.UserID = userID.String
				player.User = internal.User{
					ID:       userID.String,
					Name:     userName.String,
					Username: userUsername.String,
					UserSocials: internal.UserSocials{
						ID:       socialsID.String,
						Twitch:   twitch.String,
						Twitter:  twitter.String,
						Youtube:  youtube.String,
						Facebook: facebook.String,
					},
				}
				team.Players = append(team.Players, player)
			}

			teamFound := false
			for i := range run.Teams {
				if run.Teams[i].ID == team.ID {
					run.Teams[i].Players = append(run.Teams[i].Players, player)
					teamFound = true
					break
				}
			}
			if !teamFound {
				run.Teams = append(run.Teams, team)
			}
		}

		runMap[runID.String] = &run
	}

	for _, run := range runMap {
		runs = append(runs, *run)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) FindById(id string) (run internal.Run, err error) {
	query := "SELECT r.id AS run_id, r.name AS run_name, r.start_time_mili, r.estimate_string, r.estimate_mili, rm.category, rm.platform, rm.twitch_game_name, rm.note, r.schedule_id, t.id AS team_id, t.name AS team_name, u.id AS user_id, u.name AS user_name, u.username AS user_username, um.id AS user_socials_id, um.twitch AS user_twitch, um.twitter AS user_twitter, um.youtube AS user_youtube, um.facebook AS user_facebook FROM  runs AS r JOIN  run_metadata AS rm ON r.id = rm.run_id LEFT JOIN teams_runs AS tr ON r.id = tr.run_id LEFT JOIN teams AS t ON tr.team_id = t.id LEFT JOIN players AS pl ON t.id = pl.team_id LEFT JOIN users AS u ON pl.user_id = u.id LEFT JOIN user_socials AS um ON u.id = um.user_id WHERE r.id = ?;"
	rows, err := r.db.Query(query, id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer rows.Close()
	var teams []internal.RunTeams

	for rows.Next() {
		var team internal.RunTeams
		var player internal.RunTeamPlayers
		var teamID, userID sql.NullString
		var teamName, userName, userUsername, socialsID, twitch, twitter, youtube, facebook sql.NullString

		err = rows.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.Note, &run.ScheduleId, &teamID, &teamName, &userID, &userName, &userUsername, &socialsID, &twitch, &twitter, &youtube, &facebook)
		if err != nil {
			if err == sql.ErrNoRows {
				err = internal.ErrRunRepositoryNotFound
			}
			logger.Log.Error(err.Error())
			return
		}

		if teamID.Valid {
			team.ID = teamID.String
			team.Name = teamName.String

			if userID.Valid {
				player.UserID = userID.String
				player.User = internal.User{
					ID:       userID.String,
					Name:     userName.String,
					Username: userUsername.String,
					UserSocials: internal.UserSocials{
						ID:       socialsID.String,
						Twitch:   twitch.String,
						Twitter:  twitter.String,
						Youtube:  youtube.String,
						Facebook: facebook.String,
					},
				}
				team.Players = append(team.Players, player)
			}

			teamFound := false
			for i := range teams {
				if teams[i].ID == team.ID {
					teams[i].Players = append(teams[i].Players, player)
					teamFound = true
					break
				}
			}
			if !teamFound {
				teams = append(teams, team)
			}
		}
	}

	run.Teams = teams

	return
}

func (r *RunSqlite) Save(run *internal.Run) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO runs (id, name, start_time_mili, estimate_string, estimate_mili, schedule_id) VALUES (?, ?, ?, ?, ?, ?);", run.ID, run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.ScheduleId)
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

	_, err = tx.Exec("INSERT INTO run_metadata (id, run_id, category, platform, twitch_game_name, run_name, note) VALUES (?, ?, ?, ?, ?, ?, ?);", run.ID, run.ID, run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.RunName, run.RunMetadata.Note)
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

	for _, team := range run.Teams {
		_, err = tx.Exec("INSERT INTO teams (id, name) VALUES (?, ?) ON CONFLICT(id) DO NOTHING;", team.ID, team.Name)
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

		_, err = tx.Exec("INSERT INTO teams_runs (run_id, team_id) VALUES (?, ?);", run.ID, team.ID)
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

		for _, player := range team.Players {
			_, err = tx.Exec("INSERT INTO players (team_id, user_id) VALUES (?, ?);", team.ID, player.UserID)
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
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) Update(run *internal.Run) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	// Actualizar la entidad Run
	_, err = tx.Exec("UPDATE runs SET name = ?, start_time_mili = ?, estimate_string = ?, estimate_mili = ?, schedule_id = ? WHERE id = ?;", run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.ScheduleId, run.ID)
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

	// Actualizar los metadatos de Run
	_, err = tx.Exec("UPDATE run_metadata SET category = ?, platform = ?, twitch_game_name = ?, run_name = ?, note = ? WHERE run_id = ?;", run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.RunName, run.RunMetadata.Note, run.ID)
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

	// Eliminar asociaciones antiguas de teams_runs y players
	_, err = tx.Exec("DELETE FROM teams_runs WHERE run_id = ?;", run.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Eliminar jugadores antiguos asociados con los equipos
	for _, team := range run.Teams {
		_, err = tx.Exec("DELETE FROM players WHERE team_id = ?;", team.ID)
		if err != nil {
			logger.Log.Error("Error deleting from players: " + err.Error())
			return
		}
	}

	// Insertar o actualizar equipos y jugadores
	for _, team := range run.Teams {
		// Insertar o actualizar el equipo
		_, err = tx.Exec("INSERT INTO teams (id, name) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET name = excluded.name;", team.ID, team.Name)
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

		// Insertar asociación en teams_runs
		_, err = tx.Exec("INSERT INTO teams_runs (run_id, team_id) VALUES (?, ?);", run.ID, team.ID)
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

		// Insertar jugadores en el equipo
		for _, player := range team.Players {
			_, err = tx.Exec("INSERT INTO players (team_id, user_id) VALUES (?, ?) ON CONFLICT(team_id, user_id) DO UPDATE SET team_id = excluded.team_id;", team.ID, player.UserID)
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
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) Delete(id string) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error("Error starting transaction on delete run: " + err.Error())
		return
	}

	// Eliminar teams asociados con la run
	_, err = tx.Exec("DELETE FROM teams WHERE id IN (SELECT team_id FROM teams_runs WHERE run_id = ?);", id)
	if err != nil {
		logger.Log.Error("Error deleting from teams: " + err.Error())
		return
	}

	res, err := tx.Exec("DELETE FROM runs WHERE id = ?;", id)
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

	if err = tx.Commit(); err != nil {
		logger.Log.Error("Error committing transaction on delete run: " + err.Error())
		return
	}

	return
}
