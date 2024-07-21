package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type BidSqlite struct {
	db *sql.DB
}

func NewBidSqlite(db *sql.DB) *BidSqlite {
	return &BidSqlite{db}
}

func (r *BidSqlite) FindAll() (bids []internal.Bid, err error) {
	rows, err := r.db.Query("SELECT b.`id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description`, b.`type`, b.`create_new_options`, b.`run_id`, bo.`id`, bo.`name` AS bid_option_name, bo.`current_amount`, bo.`bid_id` AS bid_option_current_amount FROM `bids` AS b LEFT JOIN bid_options AS bo ON b.`id` = bo.`bid_id`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var bid internal.Bid
		var bidOption internal.BidOptionsSQL
		err = rows.Scan(&bid.ID, &bid.Bidname, &bid.Goal, &bid.CurrentAmount, &bid.Description, &bid.Type, &bid.CreateNewOptions, &bid.RunID, &bidOption.ID, &bidOption.Name, &bidOption.CurrentAmount, &bidOption.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		exists := false
		for i := range bids {
			if bids[i].ID == bid.ID {
				if bidOption.ID.Valid {
					bids[i].BidOptions = append(bids[i].BidOptions, internal.BidOptions{
						ID:            bidOption.ID.String,
						Name:          bidOption.Name.String,
						CurrentAmount: float64(bidOption.CurrentAmount.Float64),
						BidID:         bidOption.BidID.String,
					})
				}
				exists = true
				break
			}
		}

		if !exists {
			if bidOption.ID.Valid {
				bid.BidOptions = []internal.BidOptions{{
					ID:            bidOption.ID.String,
					Name:          bidOption.Name.String,
					CurrentAmount: float64(bidOption.CurrentAmount.Float64),
					BidID:         bidOption.BidID.String,
				}}
			}
			bids = append(bids, bid)
		}
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *BidSqlite) FindById(id string) (bid internal.Bid, err error) {
	row := r.db.QueryRow("SELECT b.`id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description`, b.`type`, b.`create_new_options`, b.`run_id` FROM `bids` AS b WHERE b.`id` = ?", id)

	err = row.Scan(&bid.ID, &bid.Bidname, &bid.Goal, &bid.CurrentAmount, &bid.Description, &bid.Type, &bid.CreateNewOptions, &bid.RunID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrBidRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	rows, err := r.db.Query("SELECT bo.`id`, bo.`name` AS bid_option_name, bo.`current_amount`, bo.`bid_id` AS bid_option_current_amount FROM `bids` AS b JOIN bid_options AS bo ON b.`id` = bo.`bid_id` WHERE b.`id` = ?;", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var bidOption internal.BidOptions
		err = rows.Scan(&bidOption.ID, &bidOption.Name, &bidOption.CurrentAmount, &bidOption.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		if bidOption.ID != "" {
			bid.BidOptions = append(bid.BidOptions, bidOption)
		}
	}
	return
}

func (r *BidSqlite) Save(bid *internal.Bid) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
	}
	_, err = tx.Exec("INSERT INTO `bids` (`id`, `bidname`, goal, current_amount, description, type, create_new_options, run_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", bid.ID, bid.Bidname, bid.Goal, bid.CurrentAmount, bid.Description, bid.Type, bid.CreateNewOptions, bid.RunID)
	if err != nil {
		tx.Rollback()
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrBidRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	query := "INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?)"

	for _, option := range bid.BidOptions {
		_, err = tx.Exec(query, option.ID, option.BidID, option.Name, option.CurrentAmount)
		if err != nil {
			tx.Rollback()
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) {
				switch sqliteErr.ExtendedCode {
				case sqlite3.ErrConstraintUnique:
					err = internal.ErrBidRepositoryDuplicated
				default:
					err = internal.ErrBidDatabase
				}
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

func (r *BidSqlite) Update(bid *internal.Bid) (err error) {

	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
	}
	logger.Log.Info("bid: ", bid)
	_, err = tx.Exec("UPDATE `bids` SET `bidname` = ?, `goal` = ?, `current_amount` = ?, `description` = ?, `type` = ?, `create_new_options` = ?, `run_id` = ? WHERE `id` = ?;", bid.Bidname, bid.Goal, bid.CurrentAmount, bid.Description, bid.Type, bid.CreateNewOptions, bid.RunID, bid.ID)
	if err != nil {
		tx.Rollback()
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrBidRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	if bid.Type != internal.Bidwar {
		deleteBidOptionsQuery := "DELETE FROM `bid_options` WHERE bid_id = ?;"

		_, err = tx.Exec(deleteBidOptionsQuery, bid.ID)
		if err != nil {
			tx.Rollback()
			logger.Log.Error(err.Error())
		}
	} else {
		query := "INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET `name` = excluded.name, `current_amount` = excluded.current_amount"

		for _, option := range bid.BidOptions {
			_, err = tx.Exec(query, option.ID, option.BidID, option.Name, option.CurrentAmount)
			if err != nil {
				tx.Rollback()
				var sqliteErr sqlite3.Error
				if errors.As(err, &sqliteErr) {
					switch sqliteErr.ExtendedCode {
					case sqlite3.ErrConstraintUnique:
						err = internal.ErrBidRepositoryDuplicated
					default:
						err = internal.ErrBidDatabase
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

func (r *BidSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `bids` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrBidRepositoryNotFound
		return
	}

	return
}
