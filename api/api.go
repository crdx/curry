package api

import (
	"fmt"
	"log"

	"github.com/crdx/curry/util"
)

func FetchRawForDay(day string, accessKey string) []byte {
	body, err := util.HttpGetJson(
		fmt.Sprintf("http://api.exchangeratesapi.io/%s", day),
		map[string]string{
			"access_key": accessKey,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	return body
}
