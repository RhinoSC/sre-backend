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
	rows, err := r.db.Query(`
			SELECT 
					s.id, s.name, s.start_time_mili, s.end_time_mili, s.setup_time_mili, s.event_id, r.id
			FROM 
					schedules AS s 
			LEFT JOIN 
					runs AS r ON r.schedule_id = s.id
			ORDER BY 
					r.start_time_mili ASC;
	`)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer rows.Close()

	scheduleMap := make(map[string]*internal.Schedule)

	for rows.Next() {
		var scheduleID, runID sql.NullString
		var scheduleName sql.NullString
		var scheduleStartTimeMili, scheduleEndTimeMili, scheduleSetupTimeMilli sql.NullInt64
		var scheduleEventID sql.NullString

		err = rows.Scan(&scheduleID, &scheduleName, &scheduleStartTimeMili, &scheduleEndTimeMili, &scheduleSetupTimeMilli, &scheduleEventID, &runID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		schedule, exists := scheduleMap[scheduleID.String]
		if !exists {
			schedule = &internal.Schedule{
				ID:              scheduleID.String,
				Name:            scheduleName.String,
				Start_time_mili: scheduleStartTimeMili.Int64,
				End_time_mili:   scheduleEndTimeMili.Int64,
				Setup_time_mili: scheduleSetupTimeMilli.Int64,
				EventID:         scheduleEventID.String,
				Runs:            []internal.Run{},
			}
			scheduleMap[scheduleID.String] = schedule
		}

		if runID.Valid {
			run := internal.Run{
				ID: runID.String,
			}
			schedule.Runs = append(schedule.Runs, run)
		}
	}

	if err = rows.Err(); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for _, schedule := range scheduleMap {
		schedules = append(schedules, *schedule)
	}

	return
}

func (r *ScheduleSqlite) FindById(id string) (schedule internal.Schedule, err error) {
	row := r.db.QueryRow("SELECT s.`id`, s.`name`, s.`start_time_mili`, s.`end_time_mili`, s.`setup_time_mili`, s.`event_id` FROM `schedules` AS `s` WHERE s.`id` == ?;", id)

	err = row.Scan(&schedule.ID, &schedule.Name, &schedule.Start_time_mili, &schedule.End_time_mili, &schedule.Setup_time_mili, &schedule.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrScheduleRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	var run internal.Run

	// Query for the runs
	rows, err := r.db.Query(`
	SELECT r.id AS run_id, r.name AS run_name, r.start_time_mili, r.estimate_string, r.estimate_mili, r.setup_time_mili, r.status,
		rm.category, rm.platform, rm.twitch_game_name, rm.twitch_game_id, rm.note, r.schedule_id,
		t.id AS team_id, t.name AS team_name,
		u.id AS user_id, u.name AS user_name, u.username AS user_username,
		um.id AS user_socials_id, um.twitch AS user_twitch, um.twitter AS user_twitter, um.youtube AS user_youtube, um.facebook AS user_facebook,
		b.id AS bid_id, b.bidname AS bid_name, b.goal AS bid_goal, b.current_amount AS bid_current_amount, b.description AS bid_description, b.type AS bid_type, b.create_new_options AS bid_create_new_options,
		bo.id AS bid_option_id, bo.name AS bid_option_name, bo.current_amount AS bid_option_current_amount
	FROM 
		runs AS r
		JOIN run_metadata AS rm ON r.id = rm.run_id
		LEFT JOIN teams AS t ON t.run_id = r.id
		LEFT JOIN players AS pl ON t.id = pl.team_id
		LEFT JOIN users AS u ON pl.user_id = u.id
		LEFT JOIN user_socials AS um ON u.id = um.user_id
		LEFT JOIN bids AS b ON r.id = b.run_id
		LEFT JOIN bid_options AS bo ON b.id = bo.bid_id
	WHERE 
		r.schedule_id = ?
	ORDER BY 
		r.start_time_mili ASC`, id)
	if err != nil {
		return schedule, err
	}
	defer rows.Close()

	var runs = make(map[string]internal.Run)

	for rows.Next() {
		var team internal.RunTeams
		var player internal.RunTeamPlayers

		var teamID, userID, bidID, bidOptionID sql.NullString
		var teamName, userName, userUsername, socialsID, twitch, twitter, youtube, facebook sql.NullString
		var bidName, bidDescription, bidType, bidOptionName sql.NullString
		var bidGoal, bidCurrentAmount, bidOptionCurrentAmount sql.NullFloat64
		var createNewOptions sql.NullBool

		err = rows.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.SetupTimeMili, &run.Status, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.TwitchGameId, &run.RunMetadata.Note, &run.ScheduleId,
			&teamID, &teamName, &userID, &userName, &userUsername, &socialsID, &twitch, &twitter, &youtube, &facebook,
			&bidID, &bidName, &bidGoal, &bidCurrentAmount, &bidDescription, &bidType, &createNewOptions, &bidOptionID, &bidOptionName, &bidOptionCurrentAmount)
		if err != nil {
			if err == sql.ErrNoRows {
				err = internal.ErrRunRepositoryNotFound
			}
			logger.Log.Error(err.Error())
			return
		}

		// Inicializar la run si no existe aún
		if _, exists := runs[run.ID]; !exists {
			run.Teams = []internal.RunTeams{}
			run.Bids = []internal.Bid{}
			runs[run.ID] = run
		}

		// Referencia a la run actual en el mapa
		currentRun := runs[run.ID]

		// Procesar equipos y jugadores
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
			for i := range currentRun.Teams {
				if currentRun.Teams[i].ID == team.ID {
					// Añadir jugadores solo si no existen
					for _, newPlayer := range team.Players {
						exists := false
						for _, existingPlayer := range currentRun.Teams[i].Players {
							if existingPlayer.UserID == newPlayer.UserID {
								exists = true
								break
							}
						}
						if !exists {
							currentRun.Teams[i].Players = append(currentRun.Teams[i].Players, newPlayer)
						}
					}
					teamFound = true
					break
				}
			}
			if !teamFound {
				currentRun.Teams = append(currentRun.Teams, team)
			}
		}

		// Procesar bids y bid_options
		if bidID.Valid {
			bidExists := false
			for i := range currentRun.Bids {
				if currentRun.Bids[i].ID == bidID.String {
					bidExists = true

					// Añadir bid_options si no existen
					if bidOptionID.Valid {
						optionExists := false
						for _, option := range currentRun.Bids[i].BidOptions {
							if option.ID == bidOptionID.String {
								optionExists = true
								break
							}
						}
						if !optionExists {
							currentRun.Bids[i].BidOptions = append(currentRun.Bids[i].BidOptions, internal.BidOptions{
								ID:            bidOptionID.String,
								Name:          bidOptionName.String,
								CurrentAmount: bidOptionCurrentAmount.Float64,
								BidID:         bidID.String,
							})
						}
					}
					break
				}
			}

			// Si el bid no existe, lo creamos
			if !bidExists {
				newBid := internal.Bid{
					ID:               bidID.String,
					Bidname:          bidName.String,
					Goal:             bidGoal.Float64,
					CurrentAmount:    bidCurrentAmount.Float64,
					Description:      bidDescription.String,
					Type:             internal.BidType(bidType.String),
					CreateNewOptions: createNewOptions.Bool,
					RunID:            run.ID,
					BidOptions:       []internal.BidOptions{},
				}

				if bidOptionID.Valid {
					newBid.BidOptions = append(newBid.BidOptions, internal.BidOptions{
						ID:            bidOptionID.String,
						Name:          bidOptionName.String,
						CurrentAmount: bidOptionCurrentAmount.Float64,
						BidID:         bidID.String,
					})
				}

				currentRun.Bids = append(currentRun.Bids, newBid)
			}
		}

		// Actualizar la run en el mapa
		runs[run.ID] = currentRun
	}

	if err = rows.Err(); err != nil {
		logger.Log.Error(err.Error())
		return
	}

	var availableRuns []internal.Run
	var orderedRuns []internal.Run
	var backupRuns []internal.Run

	for _, run := range runs {
		availableRuns = append(availableRuns, run)
		orderedRuns = append(orderedRuns, run)
		// Por ahora los backupRuns se pueden manejar igual
		backupRuns = append(backupRuns, run)
	}

	schedule.Runs = availableRuns
	schedule.OrderedRuns = orderedRuns
	schedule.BackupRuns = backupRuns

	return schedule, nil
}

func (r *ScheduleSqlite) Save(schedule *internal.Schedule) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO `schedules` (`id`, `name`, `start_time_mili`, `end_time_mili`, `setup_time_mili`, `event_id`) VALUES (?, ?, ?, ?, ?, ?)"
	_, err = tx.Exec(query, schedule.ID, schedule.Name, schedule.Start_time_mili, schedule.End_time_mili, schedule.Setup_time_mili, schedule.EventID)
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

	query = "UPDATE `events` SET `schedule_id` = ? WHERE `id` = ?;"
	_, err = tx.Exec(query, schedule.ID, schedule.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *ScheduleSqlite) Update(schedule *internal.Schedule) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	query := "UPDATE `schedules` SET `name` = ?, `start_time_mili` = ?, `end_time_mili` = ?, `setup_time_mili` = ?, `event_id` = ? WHERE `id` = ?;"
	_, err = r.db.Exec(query, schedule.Name, schedule.Start_time_mili, schedule.End_time_mili, schedule.Setup_time_mili, schedule.EventID, schedule.ID)
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

	query = "UPDATE `events` SET `schedule_id` = ? WHERE `id` = ?;"
	_, err = tx.Exec(query, schedule.ID, schedule.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error())
		return
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
