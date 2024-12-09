package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/RhinoSC/sre-backend/internal/handler/util/team_helper"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type TeamDefault struct {
	sv internal.TeamService
}

func NewTeamDefault(sv internal.TeamService) *TeamDefault {
	return &TeamDefault{
		sv: sv,
	}
}

func (h *TeamDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		teams, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Teams not found")
			return
		}

		// response

		// deserialize teams to TeamAsJSON
		data := make([]team_helper.TeamAsJSON, len(teams))
		for i, team := range teams {
			data[i] = team_helper.TeamAsJSON{
				ID:   team.ID,
				Name: team.Name,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *TeamDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		team, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrTeamServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Team not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := team_helper.TeamAsJSON{
			ID:   team.ID,
			Name: team.Name,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *TeamDefault) Create() http.HandlerFunc {
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
		var body team_helper.TeamAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid team body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		team := internal.Team{
			ID:   uuid.NewString(),
			Name: body.Name,
		}

		err = h.sv.Save(&team)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrTeamServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Team already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := team_helper.TeamAsJSON{
			ID:   team.ID,
			Name: team.Name,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *TeamDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		team, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrTeamServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Team not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		teamBody := team_helper.TeamAsJSON{
			ID:   team.ID,
			Name: team.Name,
		}

		if err := util.RequestJSON(r, &teamBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		team = internal.Team{
			ID:   teamBody.ID,
			Name: teamBody.Name,
		}

		err = h.sv.Update(&team)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrTeamServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Team already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := team_helper.TeamAsJSON{
			ID:   team.ID,
			Name: team.Name,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *TeamDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrTeamServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Team not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
