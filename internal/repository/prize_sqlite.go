package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type PrizeSqlite struct {
	db *sql.DB
}

func NewPrizeSqlite(db *sql.DB) *PrizeSqlite {
	return &PrizeSqlite{db}
}

func (r *PrizeSqlite) FindAll() (prizes []internal.Prize, err error) {
	rows, err := r.db.Query("SELECT p.`id`, p.`name`, p.`description`, p.`url`, p.`min_amount`, p.`status`, p.`international_delivery`, p.`event_id` FROM `prizes` AS `p`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var prize internal.Prize
		err = rows.Scan(&prize.ID, &prize.Name, &prize.Description, &prize.Url, &prize.MinAmount, &prize.Status, &prize.InternationalDelivery, &prize.EventID)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		prizes = append(prizes, prize)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *PrizeSqlite) FindById(id string) (prize internal.Prize, err error) {
	row := r.db.QueryRow("SELECT p.`id`, p.`name`, p.`description`, p.`url`, p.`min_amount`, p.`status`, p.`international_delivery`, p.`event_id` FROM `prizes` AS `p` WHERE p.`id` == ?;", id)

	err = row.Scan(&prize.ID, &prize.Name, &prize.Description, &prize.Url, &prize.MinAmount, &prize.Status, &prize.InternationalDelivery, &prize.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrPrizeRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *PrizeSqlite) Save(prize *internal.Prize) (err error) {
	_, err = r.db.Exec("INSERT INTO `prizes` (`id`, `name`, `description`, `url`, `min_amount`, `status`, `international_delivery`, `event_id`) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", prize.ID, prize.Name, prize.Description, prize.Url, prize.MinAmount, prize.Status, prize.InternationalDelivery, prize.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrPrizeRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *PrizeSqlite) Update(prize *internal.Prize) (err error) {
	_, err = r.db.Exec("UPDATE `prizes` SET `name` = ?, `description` = ?, `url` = ?, `min_amount` = ?, `status` = ?, `international_delivery` = ?, `event_id` = ? WHERE `id` = ?;", prize.ID, prize.Name, prize.Description, prize.Url, prize.MinAmount, prize.Status, prize.InternationalDelivery, prize.EventID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrPrizeRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *PrizeSqlite) Delete(id string) (err error) {
	res, err := r.db.Exec("DELETE FROM `prizes` WHERE `id` = ?", id)
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Log.Error(err.Error())
	}

	if rowsAffected == 0 {
		err = internal.ErrPrizeRepositoryNotFound
		return
	}

	return
}
