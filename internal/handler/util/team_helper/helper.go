package team_helper

type TeamAsJSON struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TeamAsBodyJSON struct {
	Name string `json:"name" validate:"required"`
}
