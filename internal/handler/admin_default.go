package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/auth"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/RhinoSC/sre-backend/internal/handler/util/admin_helper"
	"github.com/RhinoSC/sre-backend/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type AdminDefault struct {
	sv internal.AdminService
}

func NewAdminDefault(sv internal.AdminService) *AdminDefault {
	return &AdminDefault{
		sv: sv,
	}
}

func (h *AdminDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		admins, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Admins not found")
			return
		}

		// response

		// deserialize admins to AdminAsJSON
		data := make([]admin_helper.AdminAsJSON, len(admins))
		for i, admin := range admins {
			data[i] = admin_helper.AdminAsJSON{
				ID:       admin.ID,
				Username: admin.Username,
				Password: admin.Password,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *AdminDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		admin, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Admin not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := admin_helper.AdminAsJSON{
			ID:       admin.ID,
			Username: admin.Username,
			Password: admin.Password,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *AdminDefault) Create() http.HandlerFunc {
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
		var body admin_helper.AdminAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid admin body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		hashedPassword, err := utils.HashPassword(body.Password)
		if err != nil {
			util.ResponseError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		admin := internal.Admin{
			User: internal.User{
				ID:       uuid.NewString(),
				Username: body.Username,
			},
			Password: hashedPassword,
		}

		err = h.sv.Save(&admin)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Admin already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := admin_helper.AdminAsJSON{
			ID:       admin.ID,
			Username: admin.Username,
			Password: admin.Password,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *AdminDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		admin, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Admin not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		adminBody := admin_helper.AdminAsJSON{
			ID:       admin.ID,
			Username: admin.Username,
			Password: admin.Password,
		}

		if err := util.RequestJSON(r, &adminBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		admin = internal.Admin{
			User: internal.User{
				ID:       adminBody.ID,
				Username: adminBody.Username,
			},
			Password: adminBody.Password,
		}

		err = h.sv.Update(&admin)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Admin already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := admin_helper.AdminAsJSON{
			ID:       admin.ID,
			Username: admin.Username,
			Password: admin.Password,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *AdminDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process

		err := h.sv.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Admin not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}

func (h *AdminDefault) Login() http.HandlerFunc {
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
		var body admin_helper.AdminAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid admin body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		admin, err := h.sv.Login(body.Username, body.Password)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrAdminServiceNotFound):
				util.ResponseError(w, http.StatusConflict, "Invalid credentials")
			case errors.Is(err, internal.ErrAdminServiceInvalidPassword):
				util.ResponseError(w, http.StatusConflict, "Invalid credentials")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		token, err := auth.GenerateToken(admin.ID)
		if err != nil {
			http.Error(w, "Could not generate token", http.StatusInternalServerError)
			return
		}

		// response

		data := admin_helper.AdminAsJSON{
			ID:    admin.ID,
			Token: token,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *AdminDefault) ValidateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    "ok",
		})
	}
}
