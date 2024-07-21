package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/RhinoSC/sre-backend/internal/repository/utils/donation_helper"
	"github.com/mattn/go-sqlite3"
)

type DonationSqlite struct {
	db *sql.DB
}

func NewDonationSqlite(db *sql.DB) *DonationSqlite {
	return &DonationSqlite{db}
}

func (r *DonationSqlite) FindAll() (donations []internal.Donation, err error) {
	rows, err := r.db.Query("SELECT d.`id`, d.`name`, d.`email`, d.`time_mili`, d.`amount`, d.`description`, d.`to_bid`, d.`event_id` FROM `donations` AS d;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var donation internal.Donation
		err = rows.Scan(&donation.ID, &donation.Name, &donation.Email, &donation.TimeMili, &donation.Amount, &donation.Description, &donation.ToBid, &donation.EventID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donations = append(donations, donation)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *DonationSqlite) FindAllWithBidDetails() (donations []internal.DonationWithBidDetails, err error) {
	rows, err := r.db.Query("SELECT d.`id` AS `donation_id`, d.`name` AS `donation_name`, d.`email`, d.`time_mili`, d.`amount`, d.`description` AS `donation_description`, d.`to_bid`, d.`event_id`, b.`id` AS `bid_id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description` as `bid_description`, b.`type`, b.`create_new_options`, b.`run_id`, bo.`id` AS `option_id`, bo.`name` AS `option_name`, bo.`current_amount` AS `option_amount` FROM `donations` AS `d` LEFT JOIN `donation_bids` AS db ON d.id = db.`donation_id` LEFT JOIN `bids` AS b ON db.`bid_id` = b.`id` LEFT JOIN `bid_options` AS bo ON db.`bid_option_id` = bo.`id`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
		err = rows.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationWithBidDetailsDB.BidID, &donationWithBidDetailsDB.Bidname, &donationWithBidDetailsDB.Goal, &donationWithBidDetailsDB.CurrentAmount, &donationWithBidDetailsDB.BidDescription, &donationWithBidDetailsDB.Type, &donationWithBidDetailsDB.CreateNewOptions, &donationWithBidDetailsDB.RunID, &donationWithBidDetailsDB.OptionID, &donationWithBidDetailsDB.OptionName, &donationWithBidDetailsDB.OptionAmount)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donation := donation_helper.ConvertDonationWithBidDetailsDBtoInternal(donationWithBidDetailsDB)

		donations = append(donations, donation)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *DonationSqlite) FindById(id string) (donation internal.Donation, err error) {
	row := r.db.QueryRow("SELECT d.`id`, d.`name`, d.`email`, d.`time_mili`, d.`amount`, d.`description`, d.`to_bid`, d.`event_id` FROM `donations` AS `d` WHERE d.`id` == ?;", id)

	err = row.Scan(&donation.ID, &donation.Name, &donation.Email, &donation.TimeMili, &donation.Amount, &donation.Description, &donation.ToBid, &donation.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrDonationRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *DonationSqlite) FindByIdWithBidDetails(id string) (donation internal.DonationWithBidDetails, err error) {

	row := r.db.QueryRow("SELECT d.`id` AS `donation_id`, d.`name` AS `donation_name`, d.`email`, d.`time_mili`, d.`amount`, d.`description` AS `donation_description`, d.`to_bid`, d.`event_id`, b.`id` AS `bid_id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description` as `bid_description`, b.`type`, b.`create_new_options`, b.`run_id`, bo.`id` AS `option_id`, bo.`name` AS `option_name`, bo.`current_amount` AS `option_amount` FROM `donations` AS `d` LEFT JOIN `donation_bids` AS db ON d.id = db.`donation_id` LEFT JOIN `bids` AS b ON db.`bid_id` = b.`id` LEFT JOIN `bid_options` AS bo ON db.`bid_option_id` = bo.`id` WHERE d.`id` == ?;", id)

	var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
	err = row.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationWithBidDetailsDB.BidID, &donationWithBidDetailsDB.Bidname, &donationWithBidDetailsDB.Goal, &donationWithBidDetailsDB.CurrentAmount, &donationWithBidDetailsDB.BidDescription, &donationWithBidDetailsDB.Type, &donationWithBidDetailsDB.CreateNewOptions, &donationWithBidDetailsDB.RunID, &donationWithBidDetailsDB.OptionID, &donationWithBidDetailsDB.OptionName, &donationWithBidDetailsDB.OptionAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrDonationRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	donation = donation_helper.ConvertDonationWithBidDetailsDBtoInternal(donationWithBidDetailsDB)
	return
}

func (r *DonationSqlite) FindByEventID(id string) (donations []internal.Donation, err error) {
	rows, err := r.db.Query("SELECT d.`id`, d.`name`, d.`email`, d.`time_mili`, d.`amount`, d.`description`, d.`to_bid`, d.`event_id` FROM `donations` AS `d` WHERE d.`event_id` == ?;", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var donation internal.Donation
		err = rows.Scan(&donation.ID, &donation.Name, &donation.Email, &donation.TimeMili, &donation.Amount, &donation.Description, &donation.ToBid, &donation.EventID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donations = append(donations, donation)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *DonationSqlite) FindByEventIDWithBidDetails(id string) (donations []internal.DonationWithBidDetails, err error) {
	rows, err := r.db.Query("SELECT d.`id` AS `donation_id`, d.`name` AS `donation_name`, d.`email`, d.`time_mili`, d.`amount`, d.`description` AS `donation_description`, d.`to_bid`, d.`event_id`, b.`id` AS `bid_id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description` as `bid_description`, b.`type`, b.`create_new_options`, b.`run_id`, bo.`id` AS `option_id`, bo.`name` AS `option_name`, bo.`current_amount` AS `option_amount` FROM `donations` AS `d` LEFT JOIN `donation_bids` AS db ON d.id = db.`donation_id` LEFT JOIN `bids` AS b ON db.`bid_id` = b.`id` LEFT JOIN `bid_options` AS bo ON db.`bid_option_id` = bo.`id` WHERE d.`event_id` == ?;", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
		err = rows.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationWithBidDetailsDB.BidID, &donationWithBidDetailsDB.Bidname, &donationWithBidDetailsDB.Goal, &donationWithBidDetailsDB.CurrentAmount, &donationWithBidDetailsDB.BidDescription, &donationWithBidDetailsDB.Type, &donationWithBidDetailsDB.CreateNewOptions, &donationWithBidDetailsDB.RunID, &donationWithBidDetailsDB.OptionID, &donationWithBidDetailsDB.OptionName, &donationWithBidDetailsDB.OptionAmount)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donation := donation_helper.ConvertDonationWithBidDetailsDBtoInternal(donationWithBidDetailsDB)

		donations = append(donations, donation)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *DonationSqlite) Save(donation *internal.DonationWithBidDetails) (err error) {
	tx, err := r.db.Begin()
	defer tx.Rollback()

	if err != nil {
		logger.Log.Error(err.Error())
	}
	_, err = tx.Exec("INSERT INTO `donations` (`id`, `name`, `email`, `time_mili`, `amount`, `description`, `to_bid`, `event_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?);", donation.ID, donation.Name, donation.Email, donation.TimeMili, donation.Amount, donation.Description, donation.ToBid, donation.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrDonationRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	if donation.ToBid {
		if donation.OptionID != "" {
			query := "INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET `current_amount` = `current_amount` + excluded.current_amount;"
			_, err = tx.Exec(query, donation.OptionID, donation.BidID, donation.OptionName, donation.Amount)
			if err != nil {
				var sqliteErr sqlite3.Error
				if errors.As(err, &sqliteErr) {
					switch sqliteErr.ExtendedCode {
					case sqlite3.ErrConstraintUnique:
						err = internal.ErrDonationRepositoryDuplicated
					default:
						return
					}
					logger.Log.Error(err.Error())
					return
				}
			}
		}

		query := "INSERT INTO donation_bids (donation_id, bid_id, bid_option_id) VALUES (?, ?, ?);"
		var optionID any
		if donation.OptionID != "" {
			optionID = donation.OptionID
		} else {
			optionID = nil
		}
		_, err = tx.Exec(query, donation.ID, donation.BidID, optionID)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) {
				switch sqliteErr.ExtendedCode {
				case sqlite3.ErrConstraintUnique:
					err = internal.ErrDonationRepositoryDuplicated
				default:
					return
				}
				logger.Log.Error(err.Error())
				return
			}
		}

		query = "UPDATE bids SET current_amount = MAX((SELECT COALESCE(SUM(bid_options.current_amount), 0) FROM bid_options WHERE bid_options.bid_id = bids.`id`), ?) WHERE id = ?;"
		_, err = tx.Exec(query, donation.CurrentAmount, donation.BidID)
		if err != nil {
			var sqliteErr sqlite3.Error
			if errors.As(err, &sqliteErr) {
				switch sqliteErr.ExtendedCode {
				case sqlite3.ErrConstraintUnique:
					err = internal.ErrDonationRepositoryDuplicated
				default:
					return
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

func (r *DonationSqlite) Update(donation *internal.DonationWithBidDetails) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	// Obtener el monto original de la donación
	var originalAmount float64
	err = tx.QueryRow("SELECT amount FROM donations WHERE id = ?;", donation.ID).Scan(&originalAmount)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Calcular la diferencia entre el monto original y el nuevo monto
	amountDifference := donation.Amount - originalAmount

	// Actualizar la tabla donations
	_, err = tx.Exec("UPDATE `donations` SET `name` = ?, `email` = ?, `time_mili` = ?, `amount` = ?, `description` = ?, `to_bid` = ?, `event_id` = ? WHERE `id` = ?;",
		donation.Name, donation.Email, donation.TimeMili, donation.Amount, donation.Description, donation.ToBid, donation.EventID, donation.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	if donation.ToBid {
		if donation.OptionID != "" {
			// Actualizar la tabla bid_options
			query := "INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET `current_amount` = MAX(`current_amount` + ?, 0);"
			_, err = tx.Exec(query, donation.OptionID, donation.BidID, donation.OptionName, donation.Amount, amountDifference)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}

		// Actualizar la tabla bids

		query := "UPDATE bids SET current_amount = MAX((SELECT COALESCE(SUM(bid_options.current_amount), 0) FROM bid_options WHERE bid_options.bid_id = bids.`id`), MAX(`current_amount` + ?, 0)) WHERE id = ?;"
		_, err = tx.Exec(query, amountDifference, donation.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
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

func (r *DonationSqlite) Delete(id string) (err error) {
	tx, err := r.db.Begin()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}
	defer tx.Rollback()

	// Obtener información de la donación antes de eliminarla
	var donation internal.DonationWithBidDetails
	row := tx.QueryRow("SELECT d.`id` AS `donation_id`, d.`name` AS `donation_name`, d.`email`, d.`time_mili`, d.`amount`, d.`description` AS `donation_description`, d.`to_bid`, d.`event_id`, b.`id` AS `bid_id`, b.`bidname`, b.`goal`, b.`current_amount`, b.`description` as `bid_description`, b.`type`, b.`create_new_options`, b.`run_id`, bo.`id` AS `option_id`, bo.`name` AS `option_name`, bo.`current_amount` AS `option_amount` FROM `donations` AS `d` LEFT JOIN `donation_bids` AS db ON d.id = db.`donation_id` LEFT JOIN `bids` AS b ON db.`bid_id` = b.`id` LEFT JOIN `bid_options` AS bo ON db.`bid_option_id` = bo.`id` WHERE d.`id` == ?;", id)

	var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
	err = row.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationWithBidDetailsDB.BidID, &donationWithBidDetailsDB.Bidname, &donationWithBidDetailsDB.Goal, &donationWithBidDetailsDB.CurrentAmount, &donationWithBidDetailsDB.BidDescription, &donationWithBidDetailsDB.Type, &donationWithBidDetailsDB.CreateNewOptions, &donationWithBidDetailsDB.RunID, &donationWithBidDetailsDB.OptionID, &donationWithBidDetailsDB.OptionName, &donationWithBidDetailsDB.OptionAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrDonationRepositoryNotFound
			return
		}
		logger.Log.Error(err.Error())
		return
	}

	donation = donation_helper.ConvertDonationWithBidDetailsDBtoInternal(donationWithBidDetailsDB)

	if donation.ToBid {
		if donation.OptionID != "" {
			// Restar el monto de la opción de bid correspondiente
			_, err = tx.Exec("UPDATE `bid_options` SET `current_amount` =  MAX(`current_amount` - ?, 0) WHERE `id` = ?", donation.Amount, donation.OptionID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}

		// Restar el monto del bid correspondiente
		_, err = tx.Exec("UPDATE `bids` SET `current_amount` = MAX(`current_amount` - ?, 0) WHERE `id` = ?", donation.Amount, donation.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}
	}

	// Eliminar la donación
	res, err := tx.Exec("DELETE FROM `donations` WHERE `id` = ?", id)
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
		err = internal.ErrDonationRepositoryNotFound
		return
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}
