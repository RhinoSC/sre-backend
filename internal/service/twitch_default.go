package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/RhinoSC/sre-backend/internal"
)

var instance *TwitchDefault

type TwitchDefault struct {
	internal.Twitch
}

func CreateFirstTime(tw *internal.Twitch) *TwitchDefault {
	instance = &TwitchDefault{
		Twitch: *tw,
	}
	return instance
}

func GetTwitchInstance() *TwitchDefault {
	return instance
}

func (t *TwitchDefault) GetToken() (token string, err error) {
	url := "https://id.twitch.tv/oauth2/token"
	method := "POST"

	data := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=client_credentials", t.ClientID, t.ClientSecret)
	payload := strings.NewReader(data)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var tokenResponse internal.TwitchTokenResponse

	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	t.ClientToken = tokenResponse.AccessToken
	token = tokenResponse.AccessToken
	return
}

func (t *TwitchDefault) FindCategories(name string) (categories []internal.TwitchCategory, err error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/search/categories?query=%s", url.QueryEscape(name))
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Client-Id", t.ClientID)
	req.Header.Add("Authorization", "Bearer "+t.ClientToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		_, err = t.GetToken()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		// Retry the request with the new token
		return t.FindCategories(name)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var categoriesResponse internal.TwitchCategoriesResponse
	err = json.NewDecoder(res.Body).Decode(&categoriesResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
	categories = categoriesResponse.Data
	return
}

func (t *TwitchDefault) FindCategoryById(id int64) (category []internal.TwitchCategoryByID, err error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/games?id=%d", id)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Add("Client-Id", t.ClientID)
	req.Header.Add("Authorization", "Bearer "+t.ClientToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		_, err = t.GetToken()
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %v", err)
		}
		// Retry the request with the new token
		return t.FindCategoryById(id)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var categoryResponse internal.TwitchCategoryResponse
	err = json.NewDecoder(res.Body).Decode(&categoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(categoryResponse.Data) == 0 {
		return nil, fmt.Errorf("no category found for ID %d", id)
	}
	category = categoryResponse.Data

	return
}
