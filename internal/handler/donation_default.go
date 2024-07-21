package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/RhinoSC/sre-backend/internal/handler/util/donation_helper"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type DonationDefault struct {
	sv internal.DonationService
}

func NewDonationDefault(sv internal.DonationService) *DonationDefault {
	return &DonationDefault{
		sv: sv,
	}
}

func (h *DonationDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		withBidDetailsParam := r.URL.Query().Get("details")
		// process

		var err error
		var donations []internal.Donation
		var donationsWithBidDetails []internal.DonationWithBidDetails
		if withBidDetailsParam == "" {
			donations, err = h.sv.FindAll()
		} else {
			donationsWithBidDetails, err = h.sv.FindAllWithBidDetails()
		}

		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Donations not found")
			return
		}

		// response

		// deserialize donations to DonationAsJSON
		var data []donation_helper.DonationAsJSON
		if withBidDetailsParam == "" {
			data = donation_helper.ConvertDonationsToJSON(donations)
		} else {
			data = donation_helper.ConvertDonationsWithBidDetailsToJSON(donationsWithBidDetails)
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *DonationDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		withBidDetailsParam := r.URL.Query().Get("details")
		// process

		var err error
		var donation internal.DonationWithBidDetails
		if withBidDetailsParam == "" {
			donation.Donation, err = h.sv.FindById(id)
		} else {
			donation, err = h.sv.FindByIdWithBidDetails(id)
		}
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrDonationServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Donation not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := donation_helper.ConvertDonationWithBidDetailsToJSON(donation)

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *DonationDefault) GetByEventID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		withBidDetailsParam := r.URL.Query().Get("details")
		// process

		var err error
		var donations []internal.Donation
		var donationsWithBidDetails []internal.DonationWithBidDetails
		if withBidDetailsParam == "" {
			donations, err = h.sv.FindByEventID(id)
		} else {
			donationsWithBidDetails, err = h.sv.FindByEventIDWithBidDetails(id)
		}
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Donations not found")
			return
		}

		// response
		var data []donation_helper.DonationAsJSON
		if withBidDetailsParam == "" {
			data = donation_helper.ConvertDonationsToJSON(donations)
		} else {
			data = donation_helper.ConvertDonationsWithBidDetailsToJSON(donationsWithBidDetails)
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *DonationDefault) Create() http.HandlerFunc {
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
		var body donation_helper.DonationAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid donation body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			logger.Log.Info(err)
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		// Crear nuevo bidOption
		if body.BidDetails.CreateNewOptions && body.BidDetails.OptionID == "" && body.BidDetails.OptionName != "" {
			body.BidDetails.OptionID = uuid.NewString()
		}
		donation := internal.DonationWithBidDetails{
			Donation: internal.Donation{
				ID:          uuid.NewString(),
				Name:        body.Name,
				Email:       body.Email,
				TimeMili:    body.TimeMili,
				Amount:      body.Amount,
				Description: body.Description,
				ToBid:       *body.ToBid,
				EventID:     body.EventID,
			},
			BidID:      body.BidDetails.BidID,
			OptionID:   body.BidDetails.OptionID,
			OptionName: body.BidDetails.OptionName,
		}

		err = h.sv.Save(&donation)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrDonationServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Donation already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := donation_helper.ConvertDonationWithBidDetailsToJSON(donation)

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *DonationDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid body")
			return
		}

		var body donation_helper.DonationAsBodyJSON
		err = json.Unmarshal(requestBody, &body)
		if err != nil {
			util.ResponseError(w, http.StatusUnprocessableEntity, "Invalid donation body")
			return
		}

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			logger.Log.Info(err)
			util.ResponseError(w, http.StatusUnprocessableEntity, "Validation failed")
			return
		}

		// process

		donation, err := h.sv.FindByIdWithBidDetails(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrDonationServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Donation not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// Update donation with the new values from body
		donation.Name = body.Name
		donation.Email = body.Email
		donation.TimeMili = body.TimeMili
		donation.Amount = body.Amount
		donation.Description = body.Description
		donation.ToBid = *body.ToBid
		donation.EventID = body.EventID
		donation.BidID = body.BidDetails.BidID
		donation.OptionID = body.BidDetails.OptionID
		donation.OptionName = body.BidDetails.OptionName

		err = h.sv.Update(&donation)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrDonationServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Donation already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := donation_helper.DonationAsJSON{
			ID:          donation.ID,
			Name:        donation.Name,
			Email:       donation.Email,
			TimeMili:    donation.TimeMili,
			Amount:      donation.Amount,
			Description: donation.Description,
			ToBid:       &donation.ToBid,
			EventID:     donation.EventID,
			BidDetails: donation_helper.BidDetails{
				BidID:            donation.BidID,
				Bidname:          donation.Bidname,
				Goal:             donation.Goal,
				CurrentAmount:    donation.CurrentAmount,
				BidDescription:   donation.BidDescription,
				Type:             donation.Type,
				CreateNewOptions: donation.CreateNewOptions,
				RunID:            donation.RunID,
				OptionID:         donation.OptionID,
				OptionName:       donation.OptionName,
				OptionAmount:     donation.OptionAmount,
			},
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *DonationDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrDonationServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Donation not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, map[string]any{})
	}
}
