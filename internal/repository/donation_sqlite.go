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
		var donationBidDetails internal.DonationBidDetailsDB
		var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
		err = rows.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationBidDetails.BidID, &donationBidDetails.Bidname, &donationBidDetails.Goal, &donationBidDetails.CurrentAmount, &donationBidDetails.BidDescription, &donationBidDetails.Type, &donationBidDetails.CreateNewOptions, &donationBidDetails.RunID, &donationBidDetails.OptionID, &donationBidDetails.OptionName, &donationBidDetails.OptionAmount)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donationWithBidDetailsDB.BidDetails = &donationBidDetails
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

	var donationBidDetails internal.DonationBidDetailsDB
	var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
	err = row.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationBidDetails.BidID, &donationBidDetails.Bidname, &donationBidDetails.Goal, &donationBidDetails.CurrentAmount, &donationBidDetails.BidDescription, &donationBidDetails.Type, &donationBidDetails.CreateNewOptions, &donationBidDetails.RunID, &donationBidDetails.OptionID, &donationBidDetails.OptionName, &donationBidDetails.OptionAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrDonationRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}

	donationWithBidDetailsDB.BidDetails = &donationBidDetails
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
		var donationBidDetails internal.DonationBidDetailsDB
		var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
		err = rows.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationBidDetails.BidID, &donationBidDetails.Bidname, &donationBidDetails.Goal, &donationBidDetails.CurrentAmount, &donationBidDetails.BidDescription, &donationBidDetails.Type, &donationBidDetails.CreateNewOptions, &donationBidDetails.RunID, &donationBidDetails.OptionID, &donationBidDetails.OptionName, &donationBidDetails.OptionAmount)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		donationWithBidDetailsDB.BidDetails = &donationBidDetails
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

func (r *DonationSqlite) FindTotalDonatedByEventID(id string) (totalAmount float64, err error) {

	err = r.db.QueryRow("SELECT SUM(d.amount) FROM `donations` AS `d` WHERE d.`event_id` = ?;", id).Scan(&totalAmount)
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
		if donation.BidDetails.OptionID != "" {
			query := "INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET `current_amount` = `current_amount` + excluded.current_amount;"
			_, err = tx.Exec(query, donation.BidDetails.OptionID, donation.BidDetails.BidID, donation.BidDetails.OptionName, donation.Amount)
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
		if donation.BidDetails.OptionID != "" {
			optionID = donation.BidDetails.OptionID
		} else {
			optionID = nil
		}
		_, err = tx.Exec(query, donation.ID, donation.BidDetails.BidID, optionID)
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

		query = "UPDATE bids SET current_amount = current_amount + ? WHERE id = ?;"
		_, err = tx.Exec(query, donation.Amount, donation.BidDetails.BidID)
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

	// Actualizar la tabla donations
	_, err = tx.Exec("UPDATE `donations` SET `name` = ?, `email` = ?, `time_mili` = ?, `amount` = ?, `description` = ?, `to_bid` = ?, `event_id` = ? WHERE `id` = ?;",
		donation.Name, donation.Email, donation.TimeMili, donation.Amount, donation.Description, donation.ToBid, donation.EventID, donation.ID)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	// Si el BidDetails y NewBidDetails son los mismos
	if donation.BidDetails.BidID == donation.NewBidDetails.BidID && donation.BidDetails.OptionID == donation.NewBidDetails.OptionID {
		// Actualizar el monto en el mismo bid o bidOption
		if donation.NewBidDetails.Type == "bidwar" && donation.NewBidDetails.OptionID != "" {
			// Si es una bidwar, actualizar el bidOption
			_, err = tx.Exec("UPDATE bid_options SET current_amount = current_amount + ? WHERE id = ? AND bid_id = ?", donation.Amount-originalAmount, donation.NewBidDetails.OptionID, donation.NewBidDetails.BidID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		} else {
			// Si no es una bidwar, actualizar el bid directamente
			_, err = tx.Exec("UPDATE bids SET current_amount = current_amount + ? WHERE id = ?", donation.Amount-originalAmount, donation.NewBidDetails.BidID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}
	} else {
		// Verificar si el oldBid (BidDetails) es bidwar y manejar bidOptions
		if donation.BidDetails.Type == "bidwar" && donation.BidDetails.OptionID != "" {
			// Restar el monto de la opción en bid_options
			_, err = tx.Exec("UPDATE bid_options SET current_amount = current_amount - ? WHERE id = ? AND bid_id = ?", originalAmount, donation.BidDetails.OptionID, donation.BidDetails.BidID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}

		// Restar el monto del oldBid (BidDetails)
		_, err = tx.Exec("UPDATE bids SET current_amount = current_amount - ? WHERE id = ?", originalAmount, donation.BidDetails.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Delete en la tabla intermedia donations_bids
		_, err = tx.Exec("DELETE FROM donation_bids WHERE donation_id = ? AND bid_id = ?;", donation.ID, donation.BidDetails.BidID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		// Verificar si el newBid (NewBidDetails) es bidwar y manejar bidOptions
		if donation.NewBidDetails.Type == "bidwar" && donation.NewBidDetails.OptionID != "" {
			// Sumar el monto a la nueva opción en bid_options
			_, err = tx.Exec("INSERT INTO `bid_options` (`id`, `bid_id`, `name`, `current_amount`) VALUES (?, ?, ?, ?) ON CONFLICT(id) DO UPDATE SET `current_amount` = MAX(`current_amount` + ?, 0);",
				donation.NewBidDetails.OptionID, donation.NewBidDetails.BidID, donation.NewBidDetails.OptionName, donation.Amount, donation.Amount)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}

			// Insert en la tabla intermedia donations_bids cuando es bidwar
			_, err = tx.Exec("INSERT INTO donation_bids (donation_id, bid_id, bid_option_id) VALUES (?, ?, ?);", donation.ID, donation.NewBidDetails.BidID, donation.NewBidDetails.OptionID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		} else {
			// Insert en la tabla intermedia donations_bids cuando no es bidwar
			_, err = tx.Exec("INSERT INTO donation_bids (donation_id, bid_id) VALUES (?, ?);", donation.ID, donation.NewBidDetails.BidID, donation.NewBidDetails.OptionID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}

		// Sumar el monto al newBid (NewBidDetails)
		_, err = tx.Exec("UPDATE bids SET current_amount = current_amount + ? WHERE id = ?", donation.Amount, donation.NewBidDetails.BidID)
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

	var donationBidDetails internal.DonationBidDetailsDB
	var donationWithBidDetailsDB internal.DonationWithBidDetailsDB
	err = row.Scan(&donationWithBidDetailsDB.ID, &donationWithBidDetailsDB.Name, &donationWithBidDetailsDB.Email, &donationWithBidDetailsDB.TimeMili, &donationWithBidDetailsDB.Amount, &donationWithBidDetailsDB.Description, &donationWithBidDetailsDB.ToBid, &donationWithBidDetailsDB.EventID, &donationBidDetails.BidID, &donationBidDetails.Bidname, &donationBidDetails.Goal, &donationBidDetails.CurrentAmount, &donationBidDetails.BidDescription, &donationBidDetails.Type, &donationBidDetails.CreateNewOptions, &donationBidDetails.RunID, &donationBidDetails.OptionID, &donationBidDetails.OptionName, &donationBidDetails.OptionAmount)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrDonationRepositoryNotFound
			return
		}
		logger.Log.Error(err.Error())
		return
	}

	donationWithBidDetailsDB.BidDetails = &donationBidDetails
	donation = donation_helper.ConvertDonationWithBidDetailsDBtoInternal(donationWithBidDetailsDB)

	if donation.ToBid {
		if donation.BidDetails.OptionID != "" {
			// Restar el monto de la opción de bid correspondiente
			_, err = tx.Exec("UPDATE `bid_options` SET `current_amount` =  MAX(`current_amount` - ?, 0) WHERE `id` = ?", donation.Amount, donation.BidDetails.OptionID)
			if err != nil {
				logger.Log.Error(err.Error())
				return
			}
		}

		// Restar el monto del bid correspondiente
		_, err = tx.Exec("UPDATE `bids` SET `current_amount` = MAX(`current_amount` - ?, 0) WHERE `id` = ?", donation.Amount, donation.BidDetails.BidID)
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
