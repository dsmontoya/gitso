package bitso

import (
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
)

const (
  API_URL = "https://api.bitso.com/v2/"
  tickerPath = "ticker"
  transactionsPath = "transactions"
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

type Transaction struct {
  Amount string
  Date string
  Price string
  Tid int
  Side string
}

func TickerData(book string) (*Ticker, error) {
  ticker := &Ticker{}
  v := &url.Values{}
  v.Set("book", book)
  err := get(tickerPath, v, ticker)
  if err != nil {
    return nil, err
  }
  return ticker, nil
}

func get(path string, query *url.Values, schema interface{}) (error) {
  u, err := url.Parse(API_URL + path)
  if err != nil {
    return err
  }
  if query != nil {
    u.RawQuery = query.Encode()
  }
  resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return err
  }
  err = json.Unmarshal(body, schema)
	if err != nil {
		return err
	}
  return nil
}
