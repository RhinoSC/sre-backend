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
		r.id AS run_id, r.name AS run_name, r.start_time_mili, r.estimate_string, r.estimate_mili, r.setup_time_mili, r.status,
		rm.category, rm.platform, rm.twitch_game_name, rm.twitch_game_id, rm.note, r.schedule_id,
		t.id AS team_id, t.name AS team_name,
		u.id AS user_id, u.name AS user_name, u.username AS user_username,
		um.id AS user_socials_id, um.twitch AS user_twitch, um.twitter AS user_twitter, um.youtube AS user_youtube, um.facebook AS user_facebook,
		b.id AS bid_id, b.bidname AS bid_name, b.goal AS bid_goal, b.current_amount AS bid_current_amount, b.description AS bid_description, b.type AS bid_type, b.create_new_options AS bid_create_new_options, b.status AS bid_status,
		bo.id AS bid_option_id, bo.name AS bid_option_name, bo.current_amount AS bid_option_current_amount
	FROM 
		runs AS r
		JOIN run_metadata AS rm ON r.id = rm.run_id
		LEFT JOIN teams AS t ON t.run_id = r.id
		LEFT JOIN players AS pl ON t.id = pl.team_id
		LEFT JOIN users AS u ON pl.user_id = u.id
		LEFT JOIN user_socials AS um ON u.id = um.user_id
		LEFT JOIN bids AS b ON r.id = b.run_id
		LEFT JOIN bid_options AS bo ON b.id = bo.bid_id;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer rows.Close()

	runMap := make(map[string]*internal.Run)
	bidMap := make(map[string]*internal.Bid)

	for rows.Next() {
		var runID, teamID, userID, bidID, bidOptionID, runStatus, twitchGameId sql.NullString
		var runName, teamName, userName, userUsername, socialsID, twitch, twitter, youtube, facebook sql.NullString
		var startTimeMili, estimateMili, setupTimeMili sql.NullInt64
		var estimateString, category, platform, twitchGameName, note, scheduleID sql.NullString
		var bidName, bidDescription, bidType, bidOptionName, bidStatus sql.NullString
		var bidGoal, bidCurrentAmount, bidOptionCurrentAmount sql.NullFloat64
		var createNewOptions sql.NullBool

		err = rows.Scan(
			&runID, &runName, &startTimeMili, &estimateString, &estimateMili, &setupTimeMili, &runStatus, &category, &platform, &twitchGameName, &twitchGameId, &note, &scheduleID,
			&teamID, &teamName,
			&userID, &userName, &userUsername,
			&socialsID, &twitch, &twitter, &youtube, &facebook,
			&bidID, &bidName, &bidGoal, &bidCurrentAmount, &bidDescription, &bidType, &createNewOptions, &bidStatus,
			&bidOptionID, &bidOptionName, &bidOptionCurrentAmount,
		)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		runPtr, runExists := runMap[runID.String]
		if !runExists {
			run := internal.Run{
				ID:             runID.String,
				Name:           runName.String,
				StartTimeMili:  startTimeMili.Int64,
				EstimateString: estimateString.String,
				EstimateMili:   estimateMili.Int64,
				SetupTimeMili:  setupTimeMili.Int64,
				Status:         runStatus.String,
				RunMetadata: internal.RunMetadata{
					ID:             runID.String,
					RunID:          runID.String,
					Category:       category.String,
					Platform:       platform.String,
					TwitchGameName: twitchGameName.String,
					TwitchGameId:   twitchGameId.String,
					Note:           note.String,
				},
				ScheduleId: scheduleID.String,
				Teams:      []internal.RunTeams{},
				Bids:       []internal.Bid{},
			}
			runMap[runID.String] = &run
			runPtr = &run
		}

		if teamID.Valid {
			team := internal.RunTeams{
				ID:   teamID.String,
				Name: teamName.String,
			}

			if userID.Valid {
				player := internal.RunTeamPlayers{
					UserID: userID.String,
					User: internal.User{
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
					},
				}
				team.Players = append(team.Players, player)
			}

			teamFound := false
			for i := range runPtr.Teams {
				if runPtr.Teams[i].ID == team.ID {
					runPtr.Teams[i].Players = append(runPtr.Teams[i].Players, team.Players...)
					teamFound = true
					break
				}
			}
			if !teamFound {
				runPtr.Teams = append(runPtr.Teams, team)
			}
		}

		if bidID.Valid {
			bid, bidExists := bidMap[bidID.String]
			if !bidExists {
				bid = &internal.Bid{
					ID:               bidID.String,
					Bidname:          bidName.String,
					Goal:             bidGoal.Float64,
					CurrentAmount:    bidCurrentAmount.Float64,
					Description:      bidDescription.String,
					Type:             internal.BidType(bidType.String),
					CreateNewOptions: createNewOptions.Bool,
					Status:           bidStatus.String,
					RunID:            runID.String,
					BidOptions:       []internal.BidOptions{},
				}
				bidMap[bidID.String] = bid
				runPtr.Bids = append(runPtr.Bids, *bid)
			}

			if bidOptionID.Valid {
				option := internal.BidOptions{
					ID:            bidOptionID.String,
					Name:          bidOptionName.String,
					CurrentAmount: bidOptionCurrentAmount.Float64,
					BidID:         bidID.String,
				}
				bid.BidOptions = append(bid.BidOptions, option)
			}
		}
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
	query := `SELECT r.id AS run_id, r.name AS run_name, r.start_time_mili, r.estimate_string, r.estimate_mili, r.setup_time_mili, r.status,
						rm.category, rm.platform, rm.twitch_game_name, rm.twitch_game_id, rm.note, r.schedule_id,
						t.id AS team_id, t.name AS team_name,
						u.id AS user_id, u.name AS user_name, u.username AS user_username,
						um.id AS user_socials_id, um.twitch AS user_twitch, um.twitter AS user_twitter, um.youtube AS user_youtube, um.facebook AS user_facebook,
						b.id AS bid_id, b.bidname AS bid_name, b.goal AS bid_goal, b.current_amount AS bid_current_amount, b.description AS bid_description, b.type AS bid_type, b.create_new_options AS bid_create_new_options, b.status AS bid_status,
						bo.id AS bid_option_id, bo.name AS bid_option_name, bo.current_amount AS bid_option_current_amount FROM  runs AS r
						JOIN  run_metadata AS rm ON r.id = rm.run_id
						LEFT JOIN teams AS t ON t.run_id = r.id
						LEFT JOIN players AS pl ON t.id = pl.team_id
						LEFT JOIN users AS u ON pl.user_id = u.id
						LEFT JOIN user_socials AS um ON u.id = um.user_id
						LEFT JOIN bids AS b ON r.id = b.run_id
						LEFT JOIN bid_options AS bo ON b.id = bo.bid_id WHERE r.id = ?;`

	rows, err := r.db.Query(query, id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer rows.Close()

	var teams []internal.RunTeams
	var bids []internal.Bid
	bidMap := make(map[string]*internal.Bid)
	var runFound bool

	for rows.Next() {
		runFound = true
		var team internal.RunTeams
		var player internal.RunTeamPlayers
		var bid internal.Bid
		var bidOption internal.BidOptions

		var teamID, userID, bidID, bidOptionID sql.NullString
		var teamName, userName, userUsername, socialsID, twitch, twitter, youtube, facebook sql.NullString
		var bidName, bidDescription, bidType, bidOptionName, bidStatus sql.NullString
		var bidGoal, bidCurrentAmount, bidOptionCurrentAmount sql.NullFloat64
		var createNewOptions sql.NullBool

		err = rows.Scan(&run.ID, &run.Name, &run.StartTimeMili, &run.EstimateString, &run.EstimateMili, &run.SetupTimeMili, &run.Status, &run.RunMetadata.Category, &run.RunMetadata.Platform, &run.RunMetadata.TwitchGameName, &run.RunMetadata.TwitchGameId, &run.RunMetadata.Note, &run.ScheduleId,
			&teamID, &teamName, &userID, &userName, &userUsername, &socialsID, &twitch, &twitter, &youtube, &facebook,
			&bidID, &bidName, &bidGoal, &bidCurrentAmount, &bidDescription, &bidType, &createNewOptions, &bidStatus, &bidOptionID, &bidOptionName, &bidOptionCurrentAmount)
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

		// Procesar bids y bid_options
		if bidID.Valid {
			if _, exists := bidMap[bidID.String]; !exists {
				bid = internal.Bid{
					ID:               bidID.String,
					Bidname:          bidName.String,
					Goal:             bidGoal.Float64,
					CurrentAmount:    bidCurrentAmount.Float64,
					Description:      bidDescription.String,
					Type:             internal.BidType(bidType.String),
					CreateNewOptions: createNewOptions.Bool,
					Status:           bidStatus.String,
					RunID:            run.ID,
					BidOptions:       []internal.BidOptions{},
				}
				bidMap[bidID.String] = &bid
				bids = append(bids, bid)
			}

			if bidOptionID.Valid {
				option := internal.BidOptions{
					ID:            bidOptionID.String,
					Name:          bidOption.Name,
					CurrentAmount: bidOptionCurrentAmount.Float64,
					BidID:         bidID.String,
				}
				bidMap[bidID.String].BidOptions = append(bidMap[bidID.String].BidOptions, option)
			}
		}
	}

	if !runFound {
		err = internal.ErrRunRepositoryNotFound
		return
	}

	run.Teams = teams
	run.Bids = bids

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *RunSqlite) Save(run *internal.Run) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO runs (id, name, start_time_mili, estimate_string, estimate_mili, setup_time_mili, status, schedule_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", run.ID, run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.SetupTimeMili, run.Status, run.ScheduleId)
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

	_, err = tx.Exec("INSERT INTO run_metadata (id, run_id, category, platform, twitch_game_name, twitch_game_id, run_name, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", run.ID, run.ID, run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.TwitchGameId, run.RunMetadata.RunName, run.RunMetadata.Note)
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
		_, err = tx.Exec("INSERT INTO teams (id, name, run_id) VALUES (?, ?, ?) ON CONFLICT(id) DO NOTHING;", team.ID, team.Name, run.ID)
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

	// Insertar los bids
	for _, bid := range run.Bids {
		_, err = tx.Exec("INSERT INTO bids (id, bidname, goal, current_amount, description, type, create_new_options, status, run_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);",
			bid.ID, bid.Bidname, bid.Goal, bid.CurrentAmount, bid.Description, bid.Type, bid.CreateNewOptions, bid.Status, run.ID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Insertar las bid_options asociadas
		for _, option := range bid.BidOptions {
			_, err = tx.Exec("INSERT INTO bid_options (id, bid_id, name, current_amount) VALUES (?, ?, ?, ?);",
				option.ID, bid.ID, option.Name, option.CurrentAmount)
			if err != nil {
				logger.Log.Error(err.Error())
				return
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
	_, err = tx.Exec("UPDATE runs SET name = ?, start_time_mili = ?, estimate_string = ?, estimate_mili = ?, setup_time_mili = ?, status = ?, schedule_id = ? WHERE id = ?;", run.Name, run.StartTimeMili, run.EstimateString, run.EstimateMili, run.SetupTimeMili, run.Status, run.ScheduleId, run.ID)
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
	_, err = tx.Exec("UPDATE run_metadata SET category = ?, platform = ?, twitch_game_name = ?, twitch_game_id = ?, run_name = ?, note = ? WHERE run_id = ?;", run.RunMetadata.Category, run.RunMetadata.Platform, run.RunMetadata.TwitchGameName, run.RunMetadata.TwitchGameId, run.RunMetadata.RunName, run.RunMetadata.Note, run.ID)
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

	// Eliminar asociaciones antiguas de teams y players
	_, err = tx.Exec("DELETE FROM teams WHERE run_id = ?;", run.ID)
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
		_, err = tx.Exec("INSERT INTO teams (id, name, run_id) VALUES (?, ?) ON CONFLICT(id) DO UPDATE SET name = excluded.name;", team.ID, team.Name, run.ID)
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

	// Eliminar asociaciones antiguas de donation_bids relacionadas con el run
	_, err = tx.Exec("DELETE FROM bids WHERE run_id = ?;", run.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Insertar o actualizar bids y bid_options
	for _, bid := range run.Bids {
		_, err = tx.Exec("INSERT INTO bids (id, bidname, goal, current_amount, description, type, create_new_options, status, run_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET bidname = excluded.bidname, goal = excluded.goal, current_amount = excluded.current_amount, description = excluded.description, type = excluded.type, create_new_options = excluded.create_new_options, run_id = excluded.run_id;", bid.ID, bid.Bidname, bid.Goal, bid.CurrentAmount, bid.Description, bid.Type, bid.CreateNewOptions, bid.Status, run.ID)
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

		for _, option := range bid.BidOptions {
			// Verificar si createOptions es falso y el ID de la opción es una cadena vacía
			if !bid.CreateNewOptions && option.ID == "" {
				continue
			}
			_, err = tx.Exec("INSERT INTO bid_options (id, bid_id, name, current_amount) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET name = excluded.name, current_amount = excluded.current_amount;", option.ID, bid.ID, option.Name, option.CurrentAmount)
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

func (r *RunSqlite) UpdateRunOrder(runs []internal.Run) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("UPDATE runs SET start_time_mili = ?, status = ? WHERE id = ?")
	if err != nil {
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, run := range runs {
		_, err := stmt.Exec(run.StartTimeMili, run.Status, run.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return
}
