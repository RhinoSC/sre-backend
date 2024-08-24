package handler

import (
	"errors"
	"net/http"
	"sort"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/RhinoSC/sre-backend/internal/handler/util/run_helper"
	"github.com/RhinoSC/sre-backend/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type RunDefault struct {
	sv internal.RunService
}

func NewRunDefault(sv internal.RunService) *RunDefault {
	return &RunDefault{sv}
}

func (h *RunDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// process
		runs, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Runs not found")
			return
		}

		// response
		data := run_helper.ConvertRunsArrayToJSON(runs)

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		// process
		run, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRunServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Run not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response
		data := run_helper.ConvertRunToJSON(run)

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		var body run_helper.RunAsBodyJSON
		err := util.RequestJSON(r, &body)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// process

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		runID := uuid.NewString()
		runMetadataID := uuid.NewString()
		run := internal.Run{
			ID:             runID,
			Name:           body.Name,
			StartTimeMili:  body.StartTimeMili,
			EstimateString: body.EstimateString,
			EstimateMili:   body.EstimateMili,
			SetupTimeMili:  body.SetupTimeMili,
			Status:         body.Status,
			RunMetadata: internal.RunMetadata{
				ID:             runMetadataID,
				RunID:          runID,
				Category:       body.RunMetadata.Category,
				Platform:       body.RunMetadata.Platform,
				TwitchGameName: body.RunMetadata.TwitchGameName,
				TwitchGameId:   body.RunMetadata.TwitchGameId,
				RunName:        body.RunMetadata.RunName,
				Note:           body.RunMetadata.Note,
			},
			Teams:      make([]internal.RunTeams, len(body.RunTeams)),
			Bids:       make([]internal.Bid, len(body.Bids)),
			ScheduleId: body.ScheduleId,
		}

		for i, teamBody := range body.RunTeams {
			teamID := uuid.NewString()
			team := internal.RunTeams{
				ID:      teamID,
				Name:    teamBody.Name,
				Players: make([]internal.RunTeamPlayers, len(teamBody.Players)),
			}

			for j, playerBody := range teamBody.Players {
				player := internal.RunTeamPlayers{
					UserID: playerBody.UserID,
				}
				team.Players[j] = player
			}

			run.Teams[i] = team
		}

		// Procesar las bids
		for i, bidBody := range body.Bids {
			bidID := uuid.NewString()
			bid := internal.Bid{
				ID:               bidID,
				Bidname:          bidBody.Bidname,
				Goal:             bidBody.Goal,
				CurrentAmount:    bidBody.CurrentAmount,
				Description:      bidBody.Description,
				Type:             bidBody.Type,
				CreateNewOptions: bidBody.CreateNewOptions,
				RunID:            runID,
				BidOptions:       make([]internal.BidOptions, len(bidBody.BidOptions)),
			}

			for j, optionBody := range bidBody.BidOptions {
				bidOptionID := uuid.NewString()
				option := internal.BidOptions{
					ID:            bidOptionID,
					Name:          optionBody.Name,
					CurrentAmount: optionBody.CurrentAmount,
					BidID:         bidID,
				}
				bid.BidOptions[j] = option
			}

			run.Bids[i] = bid
		}

		err = h.sv.Save(&run)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRunServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Run already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		data := run_helper.ConvertRunToJSON(run)

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request

		id := chi.URLParam(r, "id")
		if id == "" {
			util.ResponseError(w, http.StatusBadRequest, "Invalid ID")
			return
		}

		var body run_helper.RunAsJSON
		err := util.RequestJSON(r, &body)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Process

		validate := validator.New()
		err = validate.Struct(body)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		run, err := h.sv.FindById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRunServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Run not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		run.Name = body.Name
		run.StartTimeMili = body.StartTimeMili
		run.EstimateString = body.EstimateString
		run.EstimateMili = body.EstimateMili
		run.SetupTimeMili = body.SetupTimeMili
		run.Status = body.Status
		run.RunMetadata.Category = body.RunMetadata.Category
		run.RunMetadata.Platform = body.RunMetadata.Platform
		run.RunMetadata.TwitchGameName = body.RunMetadata.TwitchGameName
		run.RunMetadata.TwitchGameId = body.RunMetadata.TwitchGameId
		run.RunMetadata.RunName = body.RunMetadata.RunName
		run.RunMetadata.Note = body.RunMetadata.Note
		run.ScheduleId = body.ScheduleId

		// Procesar equipos y jugadores
		teams := make([]internal.RunTeams, len(body.RunTeams))
		for i, teamBody := range body.RunTeams {
			teamID := run.Teams[i].ID
			if teamID == "" {
				teamID = uuid.NewString()
			}
			team := internal.RunTeams{
				ID:      teamID,
				Name:    teamBody.Name,
				Players: make([]internal.RunTeamPlayers, len(teamBody.Players)),
			}

			for j, playerBody := range teamBody.Players {
				player := internal.RunTeamPlayers{
					UserID: playerBody.UserID,
				}
				team.Players[j] = player
			}

			teams[i] = team
		}
		run.Teams = teams

		// Procesar bids y bid_options
		bids := make([]internal.Bid, len(body.Bids))
		for i, bidBody := range body.Bids {
			bidID := run.Bids[i].ID
			if bidID == "" {
				bidID = uuid.NewString()
			}
			bid := internal.Bid{
				ID:               bidID,
				Bidname:          bidBody.Bidname,
				Goal:             bidBody.Goal,
				CurrentAmount:    bidBody.CurrentAmount,
				Description:      bidBody.Description,
				Type:             bidBody.Type,
				CreateNewOptions: bidBody.CreateNewOptions,
				RunID:            run.ID,
				BidOptions:       make([]internal.BidOptions, len(bidBody.BidOptions)),
			}

			for j, optionBody := range bidBody.BidOptions {
				// Verificar si createOptions es falso y el ID de la opción es una cadena vacía
				if !bid.CreateNewOptions && optionBody.ID == "" {
					continue
				}

				bidOptionID := optionBody.ID
				if optionBody.ID == "" {
					bidOptionID = uuid.NewString()
				}

				option := internal.BidOptions{
					ID:            bidOptionID,
					Name:          optionBody.Name,
					CurrentAmount: optionBody.CurrentAmount,
				}
				bid.BidOptions[j] = option
			}

			bids[i] = bid
		}
		run.Bids = bids

		err = h.sv.Update(&run)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrRunServiceDuplicated):
				util.ResponseError(w, http.StatusConflict, "Run already exists")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// Response

		data := run_helper.ConvertRunToJSON(run)

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) Delete() http.HandlerFunc {
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
			case errors.Is(err, internal.ErrRunServiceNotFound):
				util.ResponseError(w, http.StatusNotFound, "Run not found")
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		util.ResponseJSON(w, http.StatusNoContent, nil)
	}
}

func (h *RunDefault) UpdateRunOrder() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		var body []run_helper.RunAsOrderBodyJSON
		err := util.RequestJSON(r, &body)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// process

		validate := validator.New()
		for _, bodyRun := range body {
			err = validate.Struct(bodyRun)
			if err != nil {
				logger.Log.Error(err.Error())
				util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
				return
			}
		}

		var runs []internal.Run

		for _, bodyRun := range body {
			run := internal.Run{
				ID:            bodyRun.ID,
				StartTimeMili: bodyRun.StartTimeMili,
				Status:        bodyRun.Status,
			}

			runs = append(runs, run)
		}

		err = h.sv.UpdateRunOrder(runs)
		if err != nil {
			switch {
			default:
				util.ResponseError(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}

		// response

		sort.Slice(runs, func(i, j int) bool {
			return runs[i].StartTimeMili < runs[j].StartTimeMili
		})

		var data []run_helper.RunAsOrderBodyJSON
		for _, run := range runs {
			data = append(data, run_helper.RunAsOrderBodyJSON{
				ID:            run.ID,
				StartTimeMili: run.StartTimeMili,
				Status:        run.Status,
			})
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}
