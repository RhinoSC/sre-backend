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

type PrizeAsJSON struct {
	ID                    string  `json:"id"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	Url                   string  `json:"url"`
	MinAmount             float64 `json:"min_amount"`
	Status                string  `json:"status"`
	InternationalDelivery bool    `json:"international_delivery"`
	EventID               string  `json:"event_id"`
}

type PrizeAsBodyJSON struct {
	Name                  string  `json:"name" validate:"required"`
	Description           string  `json:"description"`
	Url                   string  `json:"url" validate:"required"`
	MinAmount             float64 `json:"min_amount" validate:"required"`
	Status                string  `json:"status" validate:"required"`
	InternationalDelivery bool    `json:"international_delivery" validate:"required"`
	EventID               string  `json:"event_id" validate:"required"`
}

type PrizeDefault struct {
	sv internal.PrizeService
}

func NewPrizeDefault(sv internal.PrizeService) *PrizeDefault {
	return &PrizeDefault{
		sv: sv,
	}
}

func (h *PrizeDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		prizes, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Prizes not found")
			return
		}

		// response

		// deserialize prizes to PrizeAsJSON
		data := make([]PrizeAsJSON, len(prizes))
		for i, prize := range prizes {
			data[i] = PrizeAsJSON{
				ID:                    prize.ID,
				Name:                  prize.Name,
				Description:           prize.Description,
				Url:                   prize.Url,
				MinAmount:             prize.MinAmount,
				Status:                prize.Status,
				InternationalDelivery: prize.InternationalDelivery,
				EventID:               prize.EventID,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *PrizeDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		prize, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrPrizeServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Prize not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := PrizeAsJSON{
			ID:                    prize.ID,
			Name:                  prize.Name,
			Description:           prize.Description,
			Url:                   prize.Url,
			MinAmount:             prize.MinAmount,
			Status:                prize.Status,
			InternationalDelivery: prize.InternationalDelivery,
			EventID:               prize.EventID,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *PrizeDefault) Create() http.HandlerFunc {
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
		var body PrizeAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid prize body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		prize := internal.Prize{
			ID:                    uuid.NewString(),
			Name:                  body.Name,
			Description:           body.Description,
			Url:                   body.Url,
			MinAmount:             body.MinAmount,
			Status:                body.Status,
			InternationalDelivery: body.InternationalDelivery,
			EventID:               body.EventID,
		}

		err = h.sv.Save(&prize)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrPrizeServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Prize already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := PrizeAsJSON{
			ID:                    prize.ID,
			Name:                  prize.Name,
			Description:           prize.Description,
			Url:                   prize.Url,
			MinAmount:             prize.MinAmount,
			Status:                prize.Status,
			InternationalDelivery: prize.InternationalDelivery,
			EventID:               prize.EventID,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *PrizeDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		prize, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrPrizeServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Prize not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		prizeBody := PrizeAsJSON{
			ID:                    prize.ID,
			Name:                  prize.Name,
			Description:           prize.Description,
			Url:                   prize.Url,
			MinAmount:             prize.MinAmount,
			Status:                prize.Status,
			InternationalDelivery: prize.InternationalDelivery,
			EventID:               prize.EventID,
		}

		if err := util.RequestJSON(r, &prizeBody); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		prize = internal.Prize{
			ID:                    prizeBody.ID,
			Name:                  prizeBody.Name,
			Description:           prizeBody.Description,
			Url:                   prizeBody.Url,
			MinAmount:             prizeBody.MinAmount,
			Status:                prizeBody.Status,
			InternationalDelivery: prizeBody.InternationalDelivery,
			EventID:               prizeBody.EventID,
		}

		err = h.sv.Update(&prize)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrPrizeServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Prize already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := PrizeAsJSON{
			ID:                    prize.ID,
			Name:                  prize.Name,
			Description:           prize.Description,
			Url:                   prize.Url,
			MinAmount:             prize.MinAmount,
			Status:                prize.Status,
			InternationalDelivery: prize.InternationalDelivery,
			EventID:               prize.EventID,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *PrizeDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrPrizeServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Prize not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
