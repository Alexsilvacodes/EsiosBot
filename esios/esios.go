package esios

import (
	"encoding/json"
	"net/http"
	"time"
)

const BaseURL = "https://api.esios.ree.es/archives/70/download_json?locale=es"

var myClient = &http.Client{Timeout: 10 * time.Second}

func GetJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
