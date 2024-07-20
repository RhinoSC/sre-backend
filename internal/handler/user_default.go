package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type UserSocialsAsJSON struct {
	Twitch   string `json:"twitch" validate:"required"`
	Twitter  string `json:"twitter"`
	Youtube  string `json:"youtube"`
	Facebook string `json:"facebook"`
}

type UserAsJSON struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Username string            `json:"username"`
	Socials  UserSocialsAsJSON `json:"socials"`
}

type UserAsBodyJSON struct {
	Name     string            `json:"name" validate:"required"`
	Username string            `json:"username" validate:"required"`
	Socials  UserSocialsAsJSON `json:"socials" validate:"required"`
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
				Socials: UserSocialsAsJSON{
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
			return
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
			Socials: UserSocialsAsJSON{
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
			Socials: UserSocialsAsJSON{
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

func (h *UserDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		var mapBody map[string]any

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		if err := json.Unmarshal(requestBody, &mapBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid json")
			return
		}

		// process
		var body UserAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid user body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		user := internal.User{
			ID:       uuid.NewString(),
			Name:     body.Name,
			Username: body.Username,
			UserSocials: internal.UserSocials{
				Twitch:   body.Socials.Twitch,
				Twitter:  body.Socials.Twitter,
				Youtube:  body.Socials.Youtube,
				Facebook: body.Socials.Facebook,
			},
		}

		err = h.sv.Save(&user)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrUserServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "User already exists")
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
			Socials: UserSocialsAsJSON{
				Twitch:   user.UserSocials.Twitch,
				Twitter:  user.UserSocials.Twitter,
				Youtube:  user.UserSocials.Youtube,
				Facebook: user.UserSocials.Facebook,
			},
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *UserDefault) Update() http.HandlerFunc {
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

		userBody := UserAsJSON{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Socials: UserSocialsAsJSON{
				Twitch:   user.UserSocials.Twitch,
				Twitter:  user.UserSocials.Twitter,
				Youtube:  user.UserSocials.Youtube,
				Facebook: user.UserSocials.Facebook,
			},
		}

		if err := util.RequestJSON(r, &userBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		user = internal.User{
			ID:       userBody.ID,
			Name:     userBody.Name,
			Username: user.Username,
			UserSocials: internal.UserSocials{
				Twitch:   userBody.Socials.Twitch,
				Twitter:  userBody.Socials.Twitter,
				Youtube:  userBody.Socials.Youtube,
				Facebook: userBody.Socials.Facebook,
			},
		}

		err = h.sv.Update(&user)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrUserServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "User already exists")
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
			Socials: UserSocialsAsJSON{
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

func (h *UserDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
		}
		// process

		err := h.sv.Delete(id)
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

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
