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

type BidOptionAsJSON struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CurrentAmount float64 `json:"current_amount"`
	BidID         string  `json:"bid_id"`
}

type BidOptionAsBodyJSON struct {
	Name          string  `json:"name" validate:"required"`
	CurrentAmount float64 `json:"current_amount"`
	BidID         string  `json:"bid_id" validate:"required"`
}

type BidAsJSON struct {
	ID               string            `json:"id"`
	Bidname          string            `json:"bidname"`
	Goal             float64           `json:"goal"`
	CurrentAmount    float64           `json:"current_amount"`
	Description      string            `json:"description"`
	Type             internal.BidType  `json:"type"`
	CreateNewOptions bool              `json:"create_new_options"`
	RunID            string            `json:"run_id"`
	BidOptions       []BidOptionAsJSON `json:"bid_options"`
}

type BidAsBodyJSON struct {
	Bidname          string            `json:"bidname" validate:"required"`
	Goal             float64           `json:"goal" validate:"required"`
	CurrentAmount    float64           `json:"current_amount"`
	Description      string            `json:"description" validate:"required"`
	Type             internal.BidType  `json:"type" validate:"required"`
	CreateNewOptions bool              `json:"create_new_options"`
	RunID            string            `json:"run_id" validate:"required"`
	BidOptions       []BidOptionAsJSON `json:"bid_options"`
}

type BidDefault struct {
	sv internal.BidService
}

func NewBidDefault(sv internal.BidService) *BidDefault {
	return &BidDefault{
		sv: sv,
	}
}

func (h *BidDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		bids, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Bids not found")
			return
		}

		// response

		// deserialize bids to BidAsJSON
		data := make([]BidAsJSON, len(bids))
		for i, bid := range bids {

			var bidOptions []BidOptionAsJSON
			for _, option := range bid.BidOptions {
				optionJSON := BidOptionAsJSON{
					ID:            option.ID,
					Name:          option.Name,
					CurrentAmount: option.CurrentAmount,
					BidID:         option.BidID,
				}

				bidOptions = append(bidOptions, optionJSON)
			}
			data[i] = BidAsJSON{
				ID:               bid.ID,
				Bidname:          bid.Bidname,
				Goal:             bid.Goal,
				CurrentAmount:    bid.CurrentAmount,
				Description:      bid.Description,
				Type:             bid.Type,
				CreateNewOptions: bid.CreateNewOptions,
				RunID:            bid.RunID,
				BidOptions:       bidOptions,
			}
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BidDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}
		// process
		bid, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrBidServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Bid not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		var bidOptions []BidOptionAsJSON
		for _, option := range bid.BidOptions {
			optionJSON := BidOptionAsJSON{
				ID:            option.ID,
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         option.BidID,
			}

			bidOptions = append(bidOptions, optionJSON)
		}

		data := BidAsJSON{
			ID:               bid.ID,
			Bidname:          bid.Bidname,
			Goal:             bid.Goal,
			CurrentAmount:    bid.CurrentAmount,
			Description:      bid.Description,
			Type:             bid.Type,
			CreateNewOptions: bid.CreateNewOptions,
			RunID:            bid.RunID,
			BidOptions:       bidOptions,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BidDefault) Create() http.HandlerFunc {
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
		var body BidAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid bid body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		BidID := uuid.NewString()
		var bidOptions []internal.BidOptions
		for _, option := range body.BidOptions {
			optionJSON := internal.BidOptions{
				ID:            uuid.NewString(),
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         BidID,
			}

			bidOptions = append(bidOptions, optionJSON)
		}

		bid := internal.Bid{
			ID:               BidID,
			Bidname:          body.Bidname,
			Goal:             body.Goal,
			CurrentAmount:    body.CurrentAmount,
			Description:      body.Description,
			Type:             body.Type,
			CreateNewOptions: body.CreateNewOptions,
			RunID:            body.RunID,
			BidOptions:       bidOptions,
		}

		err = h.sv.Save(&bid)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrBidServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Bid already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		var bidOptionsResponse []BidOptionAsJSON
		for _, option := range bid.BidOptions {
			optionJSON := BidOptionAsJSON{
				ID:            option.ID,
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         option.BidID,
			}

			bidOptionsResponse = append(bidOptionsResponse, optionJSON)
		}

		data := BidAsJSON{
			ID:               bid.ID,
			Bidname:          bid.Bidname,
			Goal:             bid.Goal,
			CurrentAmount:    bid.CurrentAmount,
			Description:      bid.Description,
			Type:             bid.Type,
			CreateNewOptions: bid.CreateNewOptions,
			RunID:            bid.RunID,
			BidOptions:       bidOptionsResponse,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BidDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		bid, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrBidServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Bid not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		var bidOptionsDB []BidOptionAsJSON
		for _, option := range bid.BidOptions {
			optionJSON := BidOptionAsJSON{
				ID:            option.ID,
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         option.BidID,
			}

			bidOptionsDB = append(bidOptionsDB, optionJSON)
		}

		bidDB := BidAsJSON{
			ID:               bid.ID,
			Bidname:          bid.Bidname,
			Goal:             bid.Goal,
			CurrentAmount:    bid.CurrentAmount,
			Description:      bid.Description,
			Type:             bid.Type,
			CreateNewOptions: bid.CreateNewOptions,
			RunID:            bid.RunID,
			BidOptions:       bidOptionsDB,
		}

		if err := util.RequestJSON(r, &bidDB); err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		validate := validator.New()
		err = validate.Struct(bidDB)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		var bidOptions []internal.BidOptions

		if bidDB.Type == internal.Bidwar {
			for _, option := range bidDB.BidOptions {

				optionID := option.ID
				if optionID == "" {
					optionID = uuid.NewString()
				}

				if option.BidID == "" {
					option.BidID = bidDB.ID
				}

				option := internal.BidOptions{
					ID:            optionID,
					Name:          option.Name,
					CurrentAmount: option.CurrentAmount,
					BidID:         option.BidID,
				}

				bidOptions = append(bidOptions, option)
			}
		}

		bid = internal.Bid{
			ID:               bidDB.ID,
			Bidname:          bidDB.Bidname,
			Goal:             bidDB.Goal,
			CurrentAmount:    bidDB.CurrentAmount,
			Description:      bidDB.Description,
			Type:             bidDB.Type,
			CreateNewOptions: bidDB.CreateNewOptions,
			RunID:            bidDB.RunID,
			BidOptions:       bidOptions,
		}

		err = h.sv.Update(&bid)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrBidServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Bid already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		var bidOptionsResponse []BidOptionAsJSON
		for _, option := range bid.BidOptions {
			optionJSON := BidOptionAsJSON{
				ID:            option.ID,
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         option.BidID,
			}

			bidOptionsResponse = append(bidOptionsResponse, optionJSON)
		}

		data := BidAsJSON{
			ID:               bid.ID,
			Bidname:          bid.Bidname,
			Goal:             bid.Goal,
			CurrentAmount:    bid.CurrentAmount,
			Description:      bid.Description,
			Type:             bid.Type,
			CreateNewOptions: bid.CreateNewOptions,
			RunID:            bid.RunID,
			BidOptions:       bidOptionsResponse,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *BidDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrBidServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Bid not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
