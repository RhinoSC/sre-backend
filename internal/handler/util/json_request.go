package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrRequestContentTypeNotJSON = errors.New("request content type is not application/json")

	ErrRequestJSONInvalid = errors.New("invalid request json")
)

func RequestJSON(r *http.Request, ptr any) (err error) {
	if r.Header.Get("Content-Type") != "application/json" {
		err = ErrRequestContentTypeNotJSON
		return
	}

	err = json.NewDecoder(r.Body).Decode(ptr)
	if err != nil {
		err = fmt.Errorf("%w. %v", ErrRequestJSONInvalid, err)
		return
	}

	return
}
