package handler

import (
	"errors"
	"net/http"

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
		data := make([]run_helper.RunAsJSON, len(runs))
		for i, run := range runs {
			// Map the teams and players to the response format
			var teams []run_helper.RunTeamsAsJSON
			for _, team := range run.Teams {
				var players []run_helper.RunTeamPlayersAsJSON
				for _, player := range team.Players {

					socials := run_helper.RunTeamPlayerUserSocialsAsJSON{
						ID:       player.User.UserSocials.ID,
						Twitch:   player.User.UserSocials.Twitch,
						Twitter:  player.User.UserSocials.Twitter,
						Youtube:  player.User.UserSocials.Youtube,
						Facebook: player.User.UserSocials.Facebook,
					}
					players = append(players, run_helper.RunTeamPlayersAsJSON{
						UserID:       player.UserID,
						UserName:     player.User.Name,
						UserUsername: player.User.Username,
						Socials:      &socials,
					})
				}
				teams = append(teams, run_helper.RunTeamsAsJSON{
					ID:      team.ID,
					Name:    team.Name,
					Players: players,
				})
			}

			data[i] = run_helper.RunAsJSON{
				ID:             run.ID,
				Name:           run.Name,
				StartTimeMili:  run.StartTimeMili,
				EstimateString: run.EstimateString,
				EstimateMili:   run.EstimateMili,
				RunMetadata: run_helper.RunMetadataAsJSON{
					Category:       run.RunMetadata.Category,
					Platform:       run.RunMetadata.Platform,
					TwitchGameName: run.RunMetadata.TwitchGameName,
					RunName:        run.RunMetadata.RunName,
					Note:           run.RunMetadata.Note,
				},
				RunTeams:   teams,
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

		// Map the teams and players to the response format
		var teams []run_helper.RunTeamsAsJSON
		for _, team := range run.Teams {
			var players []run_helper.RunTeamPlayersAsJSON
			for _, player := range team.Players {

				socials := run_helper.RunTeamPlayerUserSocialsAsJSON{
					ID:       player.User.UserSocials.ID,
					Twitch:   player.User.UserSocials.Twitch,
					Twitter:  player.User.UserSocials.Twitter,
					Youtube:  player.User.UserSocials.Youtube,
					Facebook: player.User.UserSocials.Facebook,
				}

				players = append(players, run_helper.RunTeamPlayersAsJSON{
					UserID:       player.UserID,
					UserName:     player.User.Name,
					UserUsername: player.User.Username,
					Socials:      &socials,
				})
			}
			teams = append(teams, run_helper.RunTeamsAsJSON{
				ID:      team.ID,
				Name:    team.Name,
				Players: players,
			})
		}

		// response

		data := run_helper.RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: run_helper.RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			RunTeams:   teams,
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
			Teams:      make([]internal.RunTeams, len(body.RunTeams)),
		}

		for i, teamBody := range body.RunTeams {
			team := internal.RunTeams{
				ID:      uuid.NewString(),
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

		data := run_helper.RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: run_helper.RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
			RunTeams:   make([]run_helper.RunTeamsAsJSON, len(run.Teams)),
		}

		for i, team := range run.Teams {
			teamJSON := run_helper.RunTeamsAsJSON{
				ID:      team.ID,
				Name:    team.Name,
				Players: make([]run_helper.RunTeamPlayersAsJSON, len(team.Players)),
			}

			for j, player := range team.Players {
				teamJSON.Players[j] = run_helper.RunTeamPlayersAsJSON{
					UserID: player.UserID,
				}
			}

			data.RunTeams[i] = teamJSON
		}

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

		var body run_helper.RunAsBodyJSON
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
		run.RunMetadata.Category = body.RunMetadata.Category
		run.RunMetadata.Platform = body.RunMetadata.Platform
		run.RunMetadata.TwitchGameName = body.RunMetadata.TwitchGameName
		run.RunMetadata.RunName = body.RunMetadata.RunName
		run.RunMetadata.Note = body.RunMetadata.Note
		run.ScheduleId = body.ScheduleId

		logger.Log.Info(run)
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

		data := run_helper.RunAsJSON{
			ID:             run.ID,
			Name:           run.Name,
			StartTimeMili:  run.StartTimeMili,
			EstimateString: run.EstimateString,
			EstimateMili:   run.EstimateMili,
			RunMetadata: run_helper.RunMetadataAsJSON{
				Category:       run.RunMetadata.Category,
				Platform:       run.RunMetadata.Platform,
				TwitchGameName: run.RunMetadata.TwitchGameName,
				RunName:        run.RunMetadata.RunName,
				Note:           run.RunMetadata.Note,
			},
			ScheduleId: run.ScheduleId,
			RunTeams:   make([]run_helper.RunTeamsAsJSON, len(run.Teams)),
		}

		for i, team := range run.Teams {
			teamJSON := run_helper.RunTeamsAsJSON{
				ID:      team.ID,
				Name:    team.Name,
				Players: make([]run_helper.RunTeamPlayersAsJSON, len(team.Players)),
			}

			for j, player := range team.Players {

				playerJSON := run_helper.RunTeamPlayersAsJSON{
					UserID: player.UserID,
				}
				teamJSON.Players[j] = playerJSON
			}

			data.RunTeams[i] = teamJSON
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
