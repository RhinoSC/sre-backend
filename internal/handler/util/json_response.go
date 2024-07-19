package util

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ResponseJSON(w http.ResponseWriter, status int, body any) {

	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}

func ResponseError(w http.ResponseWriter, status int, message string) {
	defaultStatusCode := http.StatusInternalServerError

	if status > 299 && status < 600 {
		defaultStatusCode = status
	}

	body := errorResponse{
		Status:  http.StatusText(defaultStatusCode),
		Message: message,
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultStatusCode)
	w.Write(bytes)
}
