package internal

import "errors"

type TwitchCategory struct {
	BoxArtURL string `json:"box_art_url"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type TwitchCategoryByID struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BoxArtURL string `json:"box_art_url"`
	IgdbID    string `json:"igdb_id"`
}

type TwitchCategoryResponse struct {
	Data []TwitchCategoryByID `json:"data"`
}

type TwitchTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type TwitchCategoriesResponse struct {
	Data       []TwitchCategory `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

type Pagination struct {
	Cursor string `json:"cursor"`
}

type Twitch struct {
	ClientID     string
	ClientSecret string
	ClientToken  string
}

var (
	// SERVICE ERRORS
	ErrTwitchDatabase            = errors.New("database error")
	ErrTwitchServiceGameNotFound = errors.New("service: twitch game not found")
)

type TwitchService interface {
	GetToken() (string, error)

	FindCategories(name string) ([]TwitchCategory, error)

	FindCategoryById(id int64) ([]TwitchCategoryByID, error)
}
