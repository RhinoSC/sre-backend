package run_helper

import "github.com/RhinoSC/sre-backend/internal"

type RunTeamPlayerUserSocialsAsJSON struct {
	ID       string `json:"id,omitempty"`
	Twitch   string `json:"twitch,omitempty"`
	Twitter  string `json:"twitter,omitempty"`
	Youtube  string `json:"youtube,omitempty"`
	Facebook string `json:"facebook,omitempty"`
}

type RunTeamPlayersAsJSON struct {
	UserID       string                          `json:"id"`
	UserName     string                          `json:"name,omitempty"`
	UserUsername string                          `json:"username,omitempty"`
	Socials      *RunTeamPlayerUserSocialsAsJSON `json:"socials,omitempty"`
}

type RunTeamsAsJSON struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Players []RunTeamPlayersAsJSON `json:"players"`
}

type RunBidOptionsAsJSON struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CurrentAmount float64 `json:"current_amount"`
	BidID         string  `json:"bid_id"`
}

type RunBidsAsJSON struct {
	ID               string                `json:"id"`
	Bidname          string                `json:"bidname"`
	Goal             float64               `json:"goal"`
	CurrentAmount    float64               `json:"current_amount"`
	Description      string                `json:"description"`
	Type             internal.BidType      `json:"type"`
	CreateNewOptions bool                  `json:"create_new_options"`
	Status           string                `json:"status"`
	RunID            string                `json:"run_id"`
	BidOptions       []RunBidOptionsAsJSON `json:"bid_options"`
}

type RunMetadataAsJSON struct {
	Category       string `json:"category"`
	Platform       string `json:"platform"`
	TwitchGameName string `json:"twitch_game_name"`
	TwitchGameId   string `json:"twitch_game_id"`
	RunName        string `json:"run_name"`
	Note           string `json:"note"`
}

type RunAsJSON struct {
	ID             string             `json:"id"`
	Name           string             `json:"name,omitempty"`
	StartTimeMili  int64              `json:"start_time_mili,omitempty"`
	EstimateString string             `json:"estimate_string,omitempty"`
	EstimateMili   int64              `json:"estimate_mili,omitempty"`
	SetupTimeMili  int64              `json:"setup_time_mili,omitempty"`
	Status         string             `json:"status,omitempty"`
	RunMetadata    *RunMetadataAsJSON `json:"run_metadata,omitempty"`
	RunTeams       []RunTeamsAsJSON   `json:"teams,omitempty"`
	Bids           []RunBidsAsJSON    `json:"bids,omitempty"`
	ScheduleId     string             `json:"schedule_id,omitempty"`
}

type RunMetadataAsBodyJSON struct {
	Category       string `json:"category" validate:"required"`
	Platform       string `json:"platform" validate:"required"`
	TwitchGameName string `json:"twitch_game_name" validate:"required"`
	TwitchGameId   string `json:"twitch_game_id" validate:"required"`
	RunName        string `json:"run_name"`
	Note           string `json:"note"`
}

type RunTeamPlayersAsBodyJSON struct {
	UserID string `json:"id" validate:"required"`
}

type RunTeamsAsBodyJSON struct {
	Name    string                     `json:"name" validate:"required"`
	Players []RunTeamPlayersAsBodyJSON `json:"players" validate:"required"`
}

type RunBidOptionsAsBodyJSON struct {
	Name          string  `json:"name" validate:"required"`
	CurrentAmount float64 `json:"current_amount"`
}

type RunBidsAsBodyJSON struct {
	Bidname          string                    `json:"bidname" validate:"required"`
	Goal             *float64                  `json:"goal" validate:"required"`
	CurrentAmount    *float64                  `json:"current_amount" validate:"required"`
	Description      string                    `json:"description" validate:"required"`
	Type             internal.BidType          `json:"type" validate:"required"`
	CreateNewOptions bool                      `json:"create_new_options" validate:"required"`
	Status           string                    `json:"status" validate:"required"`
	BidOptions       []RunBidOptionsAsBodyJSON `json:"bid_options" validate:"required"`
}

type RunAsBodyJSON struct {
	Name           string                `json:"name" validate:"required"`
	StartTimeMili  int64                 `json:"start_time_mili" validate:"required"`
	EstimateString string                `json:"estimate_string" validate:"required"`
	EstimateMili   int64                 `json:"estimate_mili" validate:"required"`
	SetupTimeMili  int64                 `json:"setup_time_mili" validate:"required"`
	Status         string                `json:"status"  validate:"required"`
	RunMetadata    RunMetadataAsBodyJSON `json:"run_metadata" validate:"required"`
	RunTeams       []RunTeamsAsBodyJSON  `json:"teams" validate:"required"`
	Bids           []RunBidsAsBodyJSON   `json:"bids" validate:"required"`
	ScheduleId     string                `json:"schedule_id" validate:"required"`
}

type RunAsOrderBodyJSON struct {
	ID            string `json:"id" validate:"required"`
	StartTimeMili int64  `json:"start_time_mili" validate:"required"`
	Status        string `json:"status"  validate:"required"`
}

func ConvertRunToJSON(run internal.Run, withDetails bool) (runJSON RunAsJSON) {
	if !withDetails {
		runJSON = RunAsJSON{
			ID:   run.ID,
			Name: run.Name,
		}
		return
	}

	runJSON = RunAsJSON{
		ID:             run.ID,
		Name:           run.Name,
		StartTimeMili:  run.StartTimeMili,
		EstimateString: run.EstimateString,
		EstimateMili:   run.EstimateMili,
		SetupTimeMili:  run.SetupTimeMili,
		Status:         run.Status,
		RunMetadata: &RunMetadataAsJSON{
			Category:       run.RunMetadata.Category,
			Platform:       run.RunMetadata.Platform,
			TwitchGameName: run.RunMetadata.TwitchGameName,
			TwitchGameId:   run.RunMetadata.TwitchGameId,
			RunName:        run.RunMetadata.RunName,
			Note:           run.RunMetadata.Note,
		},
		ScheduleId: run.ScheduleId,
		RunTeams:   make([]RunTeamsAsJSON, len(run.Teams)),
		Bids:       make([]RunBidsAsJSON, len(run.Bids)),
	}

	// Llenar los equipos
	for i, team := range run.Teams {
		teamJSON := RunTeamsAsJSON{
			ID:      team.ID,
			Name:    team.Name,
			Players: make([]RunTeamPlayersAsJSON, len(team.Players)),
		}

		for j, player := range team.Players {
			teamJSON.Players[j] = RunTeamPlayersAsJSON{
				UserID:       player.UserID,
				UserName:     player.User.Name,
				UserUsername: player.User.Username,
				Socials: &RunTeamPlayerUserSocialsAsJSON{
					ID:       player.User.UserSocials.ID,
					Twitch:   player.User.UserSocials.Twitch,
					Twitter:  player.User.UserSocials.Twitter,
					Youtube:  player.User.UserSocials.Youtube,
					Facebook: player.User.UserSocials.Facebook,
				},
			}
		}

		runJSON.RunTeams[i] = teamJSON
	}

	// Llenar las bids
	for i, bid := range run.Bids {
		bidJSON := RunBidsAsJSON{
			ID:               bid.ID,
			Bidname:          bid.Bidname,
			Goal:             bid.Goal,
			CurrentAmount:    bid.CurrentAmount,
			Description:      bid.Description,
			Type:             bid.Type,
			CreateNewOptions: bid.CreateNewOptions,
			Status:           bid.Status,
			RunID:            bid.RunID,
			BidOptions:       make([]RunBidOptionsAsJSON, len(bid.BidOptions)),
		}

		for j, option := range bid.BidOptions {
			optionJSON := RunBidOptionsAsJSON{
				ID:            option.ID,
				Name:          option.Name,
				CurrentAmount: option.CurrentAmount,
				BidID:         option.BidID,
			}
			bidJSON.BidOptions[j] = optionJSON
		}

		runJSON.Bids[i] = bidJSON
	}
	return
}

func ConvertRunsArrayToJSON(runs []internal.Run, withDetails bool) (runsJSON []RunAsJSON) {
	for _, run := range runs {
		runsJSON = append(runsJSON, ConvertRunToJSON(run, withDetails))
	}

	return
}
