package bitso

import (
  "errors"
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
)

const (
  API_URL = "https://api.bitso.com/v2/"
  btcmxn = "btc_mxn"
  ethmxn = "eth_mxn"
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

// Client allows you to access to the Bitso API
type Client struct {
  configuration *Configuration
}

// Configuration stores the information needed to access
// to the private endpoints and the sanbox mode trigger.
type Configuration struct {
  Key       string
  Secret    string
  ClientId  string
  Sandbox   bool
}

// NewClient returns a new Bitso API client. It receives
// an optional Configuration, which is used to
// authenticate and enter in sandbox mode.
func NewClient(configuration *Configuration) *Client {
  return &Client{configuration}
}

// IsSandbox returns true if the sandbox mode is
// turned on.
func (c *Client) IsSandbox() bool {
  if c.configuration == nil {
    return false
  }
  return c.configuration.Sandbox
}

func (c *Client) Ticker(book string) (*Ticker, error) {
  if book != btcmxn && book != ethmxn {
    err := errors.New("Invalid book value")
    return nil, err
  }
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
