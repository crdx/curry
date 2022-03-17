package util

import (
	"io"
	"net/http"
	"time"
)

func GetYesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("2006-01-02")
}

func HttpGetJson(endpoint string, params map[string]string) (body []byte, err error) {
	client := http.Client{} // nolint
	request, err := http.NewRequest("GET", endpoint, nil)

	if err != nil {
		return
	}

	query := request.URL.Query()
	for key, value := range params {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	request.Header.Add("Accept", "application/json")
	res, err := client.Do(request)

	if err != nil {
		return
	}

	defer res.Body.Close()

	body, err = io.ReadAll(res.Body)
	return
}
