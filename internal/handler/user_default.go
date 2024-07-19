package handler

import (
	"errors"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/go-chi/chi/v5"
)

type UserAsJSON struct {
	ID       string               `json:"id"`
	Name     string               `json:"name"`
	Username string               `json:"username"`
	Socials  internal.UserSocials `json:"socials"`
}

type UserAsBodyJSON struct {
	Name     string               `json:"name"`
	Username string               `json:"username"`
	Socials  internal.UserSocials `json:"socials"`
}

type UserDefault struct {
	sv internal.UserService
}

func NewUserDefault(sv internal.UserService) *UserDefault {
	return &UserDefault{
		sv: sv,
	}
}

func (h *UserDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		users, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Users not found")
			return
		}

		// response

		// deserialize users to UserAsJSON
		data := make([]UserAsJSON, len(users))
		for i, user := range users {
			data[i] = UserAsJSON{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Socials: internal.UserSocials{
					Twitch:   user.UserSocials.Twitch,
					Twitter:  user.UserSocials.Twitter,
					Youtube:  user.UserSocials.Youtube,
					Facebook: user.UserSocials.Facebook,
				},
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *UserDefault) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
		}
		// process
		user, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrUserServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "User not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := UserAsJSON{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Socials: internal.UserSocials{
				Twitch:   user.UserSocials.Twitch,
				Twitter:  user.UserSocials.Twitter,
				Youtube:  user.UserSocials.Youtube,
				Facebook: user.UserSocials.Facebook,
			},
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *UserDefault) GetByUsername() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		username := chi.URLParam(r, "username")
		if username == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid Username")
		}
		// process
		user, err := h.sv.FindByUsername(username)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrUserServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "User not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := UserAsJSON{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Socials: internal.UserSocials{
				Twitch:   user.UserSocials.Twitch,
				Twitter:  user.UserSocials.Twitter,
				Youtube:  user.UserSocials.Youtube,
				Facebook: user.UserSocials.Facebook,
			},
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
