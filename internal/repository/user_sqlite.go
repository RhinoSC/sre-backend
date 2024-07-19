package repository

import (
	"database/sql"
	"errors"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/mattn/go-sqlite3"
)

type UserSqlite struct {
	db *sql.DB
}

func NewUserSqlite(db *sql.DB) *UserSqlite {
	return &UserSqlite{db}
}

func (r *UserSqlite) FindAll() (users []internal.User, err error) {
	rows, err := r.db.Query("SELECT u.`id`, u.`name`, u.`username`, um.`twitch`, um.`twitter`, um.`youtube`, um.`facebook` FROM `users` AS `u` JOIN user_socials AS `um` ON u.`id` = um.`user_id`;")
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	for rows.Next() {
		var user internal.User
		err = rows.Scan(&user.ID, &user.Name, &user.Username, &user.UserSocials.Twitch, &user.UserSocials.Twitter, &user.UserSocials.Youtube, &user.UserSocials.Facebook)
		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	return
}

func (r *UserSqlite) FindById(id string) (user internal.User, err error) {
	row := r.db.QueryRow("SELECT u.`id`, u.`name`, u.`username`, um.`twitch`, um.`twitter`, um.`youtube`, um.`facebook` FROM `users` AS `u` JOIN user_socials AS `um` ON u.`id` = um.`user_id` WHERE u.`id` == ?;", id)

	err = row.Scan(&user.ID, &user.Name, &user.Username, &user.UserSocials.Twitch, &user.UserSocials.Twitter, &user.UserSocials.Youtube, &user.UserSocials.Facebook)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrUserRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *UserSqlite) FindByUsername(username string) (user internal.User, err error) {
	row := r.db.QueryRow("SELECT u.`id`, u.`name`, u.`username`, um.`twitch`, um.`twitter`, um.`youtube`, um.`facebook` FROM `users` AS `u` JOIN user_socials AS `um` ON u.`id` = um.`user_id` WHERE u.`username` == ?;", username)

	err = row.Scan(&user.ID, &user.Name, &user.Username, &user.UserSocials.Twitch, &user.UserSocials.Twitter, &user.UserSocials.Youtube, &user.UserSocials.Facebook)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrUserRepositoryNotFound
		}
		logger.Log.Error(err.Error())
		return
	}
	return
}

func (r *UserSqlite) Save(user *internal.User) (err error) {
	_, err = r.db.Exec("INSERT INTO `users` (`id`, `name`, `username`) VALUES (?, ?, ?)", user.ID, user.Name, user.Username)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrUserRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	_, err = r.db.Exec("INSERT INTO `user_socials` (`id`, `user_id`, `twitch`, `twitter`, `youtube`, `facebook`) VALUES (?, ?, ?, ?, ?, ?)", user.ID, user.ID, user.UserSocials.Twitch, user.UserSocials.Twitter, user.UserSocials.Youtube, user.UserSocials.Facebook)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrUserRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	return
}

func (r *UserSqlite) Update(user *internal.User) (err error) {
	_, err = r.db.Exec("UPDATE `users` SET `name` = ?, `username` = ? WHERE `id` = ?;", user.Name, user.Username, user.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrUserRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}

	_, err = r.db.Exec("UPDATE `user_socials` SET `twitch` = ?, `twitter` = ?, `youtube` = ?, `facebook` = ? WHERE `user_id` = ?;", user.UserSocials.Twitch, user.UserSocials.Twitter, user.UserSocials.Youtube, user.UserSocials.Facebook, user.ID)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			switch sqliteErr.ExtendedCode {
			case sqlite3.ErrConstraintUnique:
				err = internal.ErrUserRepositoryDuplicated
			default:
				return
			}
			logger.Log.Error(err.Error())
			return
		}
	}
	return
}
