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

type ScheduleAsJSON struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Start_time_mili int64  `json:"start_time_mili"`
	End_time_mili   int64  `json:"end_time_mili"`
	EventID         string `json:"event_id"`
}

type ScheduleAsBodyJSON struct {
	Name            string `json:"name" validate:"required"`
	Start_time_mili int64  `json:"start_time_mili" validate:"required"`
	End_time_mili   int64  `json:"end_time_mili" validate:"required"`
	EventID         string `json:"event_id" validate:"required"`
}

type ScheduleDefault struct {
	sv internal.ScheduleService
}

func NewScheduleDefault(sv internal.ScheduleService) *ScheduleDefault {
	return &ScheduleDefault{
		sv: sv,
	}
}

func (h *ScheduleDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		schedules, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Schedules not found")
			return
		}

		// response

		// deserialize schedules to ScheduleAsJSON
		data := make([]ScheduleAsJSON, len(schedules))
		for i, schedule := range schedules {
			data[i] = ScheduleAsJSON{
				ID:              schedule.ID,
				Name:            schedule.Name,
				Start_time_mili: schedule.Start_time_mili,
				End_time_mili:   schedule.End_time_mili,
				EventID:         schedule.EventID,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *ScheduleDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		schedule, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrScheduleServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Schedule not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := ScheduleAsJSON{
			ID:              schedule.ID,
			Name:            schedule.Name,
			Start_time_mili: schedule.Start_time_mili,
			End_time_mili:   schedule.End_time_mili,
			EventID:         schedule.EventID,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *ScheduleDefault) Create() http.HandlerFunc {
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
		var body ScheduleAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid schedule body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		schedule := internal.Schedule{
			ID:              uuid.NewString(),
			Name:            body.Name,
			Start_time_mili: body.Start_time_mili,
			End_time_mili:   body.End_time_mili,
			EventID:         body.EventID,
		}

		err = h.sv.Save(&schedule)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrScheduleServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Schedule already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := ScheduleAsJSON{
			ID:              schedule.ID,
			Name:            schedule.Name,
			Start_time_mili: schedule.Start_time_mili,
			End_time_mili:   schedule.End_time_mili,
			EventID:         schedule.EventID,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *ScheduleDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		schedule, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrScheduleServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Schedule not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		scheduleBody := ScheduleAsJSON{
			ID:              schedule.ID,
			Name:            schedule.Name,
			Start_time_mili: schedule.Start_time_mili,
			End_time_mili:   schedule.End_time_mili,
			EventID:         schedule.EventID,
		}

		if err := util.RequestJSON(r, &scheduleBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		schedule = internal.Schedule{
			ID:              scheduleBody.ID,
			Name:            scheduleBody.Name,
			Start_time_mili: scheduleBody.Start_time_mili,
			End_time_mili:   scheduleBody.End_time_mili,
			EventID:         scheduleBody.EventID,
		}

		err = h.sv.Update(&schedule)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrScheduleServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Schedule already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := ScheduleAsJSON{
			ID:              schedule.ID,
			Name:            schedule.Name,
			Start_time_mili: schedule.Start_time_mili,
			End_time_mili:   schedule.End_time_mili,
			EventID:         schedule.EventID,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *ScheduleDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrScheduleServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Schedule not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
