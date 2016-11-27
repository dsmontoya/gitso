package bitso

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client allows you to access to the Bitso API
type Client struct {
	configuration *Configuration
}

// Configuration stores the information needed to access
// to the private endpoints and the sanbox mode trigger.
type Configuration struct {
	Key      string
	Secret   string
	ClientId string
	Sandbox  bool
}

type OpenOrders struct {
	fields
	Book string `json:"book,omitempty"`
}

type Order struct {
	fields
	Id     string `json:"id,omitempty"`
	Type   string `json:"type,omitempty"`
	Price  string `json:"price,omitempty"`
	Amount string `json:"amount,omitempty"`
	Status string `json:"status,omitempty"`
	Book   string `json:"book,omitempty"`
}

type Balance struct {
	fields
	MXNBalance   string `json:"mxn_balance,omitempty"`
	BTCBalance   string `json:"btc_balance,omitempty"`
	MXNReserved  string `json:"mxn_reserved,omitempty"`
	BTCReserved  string `json:"btc_reserved,omitempty"`
	MXNAvailable string `json:"mxn_available,omitempty"`
	BTCAvailable string `json:"btc_available,omitempty"`
}

// fields is included in every request made to private endpoints
type fields struct {
	Key       string `json:"key,omitempty"`
	Nonce     int64  `json:"nonce,omitempty"`
	Signature string `json:"signature,omitempty"`
	Error     *Error `json:"error,omitempty"`
}

func (a *fields) setAuthentication(key, signature string, nonce int64) {
	a.Key = key
	a.Nonce = nonce
	a.Signature = signature
}

func (a *fields) getError() *Error {
	fmt.Println("getError", a.Error)
	fmt.Println(a.Error != nil)
	if a.Error != nil {
		return a.Error
	}
	return nil
}

type Error struct {
	Code    int    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *Error) Error() string {
	fmt.Println("message", e.Message)
	fmt.Println(e.Code)
	return fmt.Sprintf("%s (code: %v)", e.Message, e.Code)
}

type authBody interface {
	setAuthentication(key, signature string, nonce int64)
	getError() *Error
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
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	ticker := &Ticker{}
	v := &url.Values{}
	v.Set("book", book)
	err := c.get(tickerPath, v, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func (c *Client) OrderBook(book string, group bool) (*OrderBook, error) {
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	orderBook := &OrderBook{}
	v := &url.Values{}
	v.Set("book", book)
	err := c.get(orderBookPath, v, orderBook)
	if err != nil {
		return nil, err
	}
	return orderBook, nil
}

/*
Transactions returns a list of recent trades from the specified book
and the specified time frame.

Valid time frames are hour and minute. Leaving time blank will set hour as the default frame.
*/
func (c *Client) Transactions(book string, time string) ([]*Transaction, error) {
	var transactions []*Transaction
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	v := &url.Values{}
	v.Set("book", book)
	v.Set("time", time)
	err := c.get(transactionsPath, v, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (c *Client) Balance() (*Balance, error) {
	balance := &Balance{}
	if err := c.post(balancePath, balance); err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Client) OpenOrders() ([]*Order, error) {
	var orders []*Order
	openOrders := &OpenOrders{}
	if err := c.post(openOrdersPath, openOrders, orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (c *Client) getSignature(nonce int64) string {
	if c.validateConfiguration() == false {
		panic("can't generate a signature without configuration")
	}
	key := c.configuration.Key
	clientId := c.configuration.ClientId
	secret := c.configuration.Secret
	message := fmt.Sprintf("%v%v%v", nonce, key, clientId)
	signature := sign(message, secret)
	return signature
}

func (c *Client) validateConfiguration() bool {
	if c.configuration == nil {
		return false
	}
	return true
}

func (c *Client) get(path string, query *url.Values, schema interface{}) error {
	if config := c.configuration; config != nil && config.Sandbox == true {
		//Mock response
	}
	u, err := url.Parse(URL + path)
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

func (c *Client) post(path string, schemas ...interface{}) error {
	var respSchema interface{}
	reqSchema := schemas[0].(authBody)
	nonce := getNonce()
	signature := c.getSignature(nonce)
	reqSchema.setAuthentication(c.configuration.Key, signature, nonce)
	payload, err := json.Marshal(reqSchema)
	if err != nil {
		return err
	}
	buff := bytes.NewBuffer(payload)
	resp, err := http.Post(URL+path, "application/json", buff)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if len(schemas) == 2 {
		respSchema = schemas[1]
	} else {
		respSchema = schemas[0]
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, respSchema)
	if err != nil {
		return err
	}
	// if inf, ok := respSchema.(authBody); ok == true {
	// 	err = inf.getError()
	// 	return err
	// }
	return nil
}
