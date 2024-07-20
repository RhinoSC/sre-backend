package handler

import (
	"errors"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
	"github.com/RhinoSC/sre-backend/internal/handler/util"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type RunMetadataAsJSON struct {
	Category       string `json:"category" validate:"required"`
	Platform       string `json:"platform" validate:"required"`
	TwitchGameName string `json:"twitch_game_name" validate:"required"`
	RunName        string `json:"run_name" validate:"required"`
	Note           string `json:"note"`
}

type RunAsJSON struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	StartTimeMili  int64             `json:"start_time_mili"`
	EstimateString string            `json:"estimate_string"`
	EstimateMili   int64             `json:"estimate_mili"`
	RunMetadata    RunMetadataAsJSON `json:"run_metadata"`
	ScheduleId     string            `json:"schedule_id"`
}

type RunAsBodyJSON struct {
	Name           string            `json:"name" validate:"required"`
	StartTimeMili  int64             `json:"start_time_mili" validate:"required"`
	EstimateString string            `json:"estimate_string" validate:"required"`
	EstimateMili   int64             `json:"estimate_mili" validate:"required"`
	RunMetadata    RunMetadataAsJSON `json:"run_metadata" validate:"required"`
	ScheduleId     string            `json:"schedule_id" validate:"required"`
}

type RunDefault struct {
	sv internal.RunService
}

func NewRunDefault(sv internal.RunService) *RunDefault {
	return &RunDefault{sv}
}

func (h *RunDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process

		runs, err := h.sv.FindAll()
		if err != nil {
			util.ResponseError(w, http.StatusNotFound, "Runs not found")
			return
		}

		// response

		data := make([]RunAsJSON, len(runs))
		for i, run := range runs {
			data[i] = RunAsJSON{
				ID:             run.ID,
				Name:           run.Name,
				StartTimeMili:  run.StartTimeMili,
				EstimateString: run.EstimateString,
				EstimateMili:   run.EstimateMili,
				RunMetadata: RunMetadataAsJSON{
					Category:       run.RunMetadata.Category,
					Platform:       run.RunMetadata.Platform,
					TwitchGameName: run.RunMetadata.TwitchGameName,
					RunName:        run.RunMetadata.RunName,
					Note:           run.RunMetadata.Note,
				},
				ScheduleId: run.ScheduleId,
			}
		}

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

		data := RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
		}

		util.ResponseJSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		var body RunAsBodyJSON
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

		run := internal.Run{
			ID:             uuid.NewString(),
			Name:           body.Name,
			StartTimeMili:  body.StartTimeMili,
			EstimateString: body.EstimateString,
			EstimateMili:   body.EstimateMili,
			RunMetadata: internal.RunMetadata{
				ID:             uuid.NewString(),
				Category:       body.RunMetadata.Category,
				Platform:       body.RunMetadata.Platform,
				TwitchGameName: body.RunMetadata.TwitchGameName,
				RunName:        body.RunMetadata.RunName,
				Note:           body.RunMetadata.Note,
			},
			ScheduleId: body.ScheduleId,
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

		data := RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
		}

		util.ResponseJSON(w, http.StatusCreated, map[string]any{
			"message": "success",
			"data":    data,
		})
	}
}

func (h *RunDefault) Update() http.HandlerFunc {
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

		runBody := RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
		}

		if err := util.RequestJSON(r, &runBody); err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		run = internal.Run{
			ID:             runBody.ID,
			Name:           runBody.Name,
			StartTimeMili:  runBody.StartTimeMili,
			EstimateString: runBody.EstimateString,
			EstimateMili:   runBody.EstimateMili,
			RunMetadata: internal.RunMetadata{
				ID:             run.RunMetadata.ID,
				Category:       runBody.RunMetadata.Category,
				Platform:       runBody.RunMetadata.Platform,
				TwitchGameName: runBody.RunMetadata.TwitchGameName,
				RunName:        runBody.RunMetadata.RunName,
				Note:           runBody.RunMetadata.Note,
			},
			ScheduleId: runBody.ScheduleId,
		}

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

		// response

		data := RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
		}

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
