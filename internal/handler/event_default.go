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

type EventAsJSON struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Start_time_mili int64  `json:"start_time_mili"`
	End_time_mili   int64  `json:"end_time_mili"`
	Schedule_id     string `json:"schedule_id"`
}

type EventAsBodyJSON struct {
	Name            string `json:"name" validate:"required"`
	Start_time_mili int64  `json:"start_time_mili" validate:"required"`
	End_time_mili   int64  `json:"end_time_mili" validate:"required"`
	Schedule_id     string `json:"schedule_id"`
}

type EventInfoCountAsJSON struct {
	Schedules_count int64 `json:"schedules"`
	Runs_count      int64 `json:"runs"`
	Prizes_count    int64 `json:"prizes"`
	Bids_count      int64 `json:"bids"`
	Donations_count int64 `json:"donations"`
	Users_count     int64 `json:"users"`
}

type EventDefault struct {
	sv internal.EventService
}

func NewEventDefault(sv internal.EventService) *EventDefault {
	return &EventDefault{
		sv: sv,
	}
}

func (h *EventDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		events, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Events not found")
			return
		}

		// response

		// deserialize events to EventAsJSON
		data := make([]EventAsJSON, len(events))
		for i, event := range events {
			data[i] = EventAsJSON{
				ID:              event.ID,
				Name:            event.Name,
				Start_time_mili: event.Start_time_mili,
				End_time_mili:   event.End_time_mili,
				Schedule_id:     event.Schedule_id,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EventDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		event, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrEventServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Event not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := EventAsJSON{
			ID:              event.ID,
			Name:            event.Name,
			Start_time_mili: event.Start_time_mili,
			End_time_mili:   event.End_time_mili,
			Schedule_id:     event.Schedule_id,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EventDefault) Create() http.HandlerFunc {
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
		var body EventAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid event body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		event := internal.Event{
			ID:              uuid.NewString(),
			Name:            body.Name,
			Start_time_mili: body.Start_time_mili,
			End_time_mili:   body.End_time_mili,
			Schedule_id:     body.Schedule_id,
		}

		err = h.sv.Save(&event)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrEventServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Event already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := EventAsJSON{
			ID:              event.ID,
			Name:            event.Name,
			Start_time_mili: event.Start_time_mili,
			End_time_mili:   event.End_time_mili,
			Schedule_id:     event.Schedule_id,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EventDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		event, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrEventServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Event not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		eventBody := EventAsJSON{
			ID:              event.ID,
			Name:            event.Name,
			Start_time_mili: event.Start_time_mili,
			End_time_mili:   event.End_time_mili,
			Schedule_id:     event.Schedule_id,
		}

		if err := util.RequestJSON(r, &eventBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		event = internal.Event{
			ID:              eventBody.ID,
			Name:            eventBody.Name,
			Start_time_mili: eventBody.Start_time_mili,
			End_time_mili:   eventBody.End_time_mili,
			Schedule_id:     eventBody.Schedule_id,
		}

		err = h.sv.Update(&event)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrEventServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Event already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := EventAsJSON{
			ID:              event.ID,
			Name:            event.Name,
			Start_time_mili: event.Start_time_mili,
			End_time_mili:   event.End_time_mili,
			Schedule_id:     event.Schedule_id,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *EventDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrEventServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Event not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}

func (h *EventDefault) GetBasicInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// process

		count, err := h.sv.GetBasicInfo()
		if err != nil {
			switch {
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		data := EventInfoCountAsJSON{
			Schedules_count: count.Schedules_count,
			Runs_count:      count.Runs_count,
			Prizes_count:    count.Prizes_count,
			Bids_count:      count.Bids_count,
			Donations_count: count.Donations_count,
			Users_count:     count.Users_count,
		}
		// response

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
