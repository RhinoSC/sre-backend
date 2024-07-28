package admin_helper

type AdminAsJSON struct {
	ID       string `json:"id"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type AdminAsBodyJSON struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
