package service

import (
	"bytes"
	"fmt"
	"net/http"
)

var layoutInstance *LayoutDefault

type LayoutDefault struct {
	url string
}

func CreateLayoutFirstTime(url string) *LayoutDefault {
	layoutInstance = &LayoutDefault{
		url: url,
	}
	return layoutInstance
}

func GetLayoutInstance() *LayoutDefault {
	return layoutInstance
}

func (l *LayoutDefault) NotifyTotalDonated(bidID string) (err error) {
	url := l.url + "/total-donated"
	method := "POST"

	data := fmt.Sprintf(`{"data":"%s", "message":"success"}`, bidID)
	body := []byte(data)
	// body := []byte(`{}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	return
}
