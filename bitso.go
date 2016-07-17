package bitso

import (
  "net/http"
  "io/ioutil"
  "encoding/json"
)

const (
  API_URL = "https://api.bitso.com/v2/"
  tickerPath = "ticker"
)

type Ticker struct {
  High string
  Last string
  Timestamp string
  Volume string
  Vwap string
  Low string
  Ask string
  Bid string
}

func get(path string, query map[string]string, schema interface{}) ([]byte, error) {
  url := API_URL + path
  resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  err = json.Unmarshal(body, schema)
	if err != nil {
		return nil, err
	}
  return body, nil
}
