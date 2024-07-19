package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RhinoSC/sre-backend/internal"
)

type UserAsJSON struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserAsBodyJSON struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type UserDefault struct {
	sv internal.UserService
}

func NewUserDefault(sv internal.UserService) *UserDefault {
	return &UserDefault{
		sv: sv,
	}
}

func (h *UserDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		users, err := h.sv.FindAll()
		if err != nil {
			body, err := json.Marshal("Users not found")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write(body)
			return
		}

		// response

		// deserialize users to UserAsJSON
		data := make([]UserAsJSON, len(users))
		for i, user := range users {
			data[i] = UserAsJSON{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
			}
		}

		body, err := json.Marshal(map[string]any{
			"message": "success",
			"data":    data,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
}
