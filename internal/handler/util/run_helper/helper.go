package run_helper

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

type RunMetadataAsJSON struct {
	Category       string `json:"category"`
	Platform       string `json:"platform"`
	TwitchGameName string `json:"twitch_game_name"`
	RunName        string `json:"run_name"`
	Note           string `json:"note"`
}

type RunAsJSON struct {
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	StartTimeMili  int64             `json:"start_time_mili"`
	EstimateString string            `json:"estimate_string"`
	EstimateMili   int64             `json:"estimate_mili"`
	RunMetadata    RunMetadataAsJSON `json:"run_metadata"`
	RunTeams       []RunTeamsAsJSON  `json:"teams,omitempty"`
	ScheduleId     string            `json:"schedule_id"`
}

type RunMetadataAsBodyJSON struct {
	Category       string `json:"category" validate:"required"`
	Platform       string `json:"platform" validate:"required"`
	TwitchGameName string `json:"twitch_game_name" validate:"required"`
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

type RunAsBodyJSON struct {
	Name           string                `json:"name" validate:"required"`
	StartTimeMili  int64                 `json:"start_time_mili" validate:"required"`
	EstimateString string                `json:"estimate_string" validate:"required"`
	EstimateMili   int64                 `json:"estimate_mili" validate:"required"`
	RunMetadata    RunMetadataAsBodyJSON `json:"run_metadata" validate:"required"`
	RunTeams       []RunTeamsAsBodyJSON  `json:"teams" validate:"required"`
	ScheduleId     string                `json:"schedule_id" validate:"required"`
}
