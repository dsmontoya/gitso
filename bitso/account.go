package bitso

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Account allows you to access to the Bitso API
type Account struct {
	keys *Keys
}

// Keys stores the information needed to access
// to the private endpoints.
type Keys struct {
	Key      string
	Secret   string
	ClientId string
}

type openOrders struct {
	fields
	Book string `json:"book,omitempty"`
}

type Order struct {
	fields
	Id       string `json:"id,omitempty"`
	Type     string `json:"type,omitempty"`
	Price    string `json:"price,omitempty"`
	Amount   string `json:"amount,omitempty"`
	Datetime string `json:"datetime,omitempty"`
	Status   string `json:"status,omitempty"`
	Book     string `json:"book,omitempty"`
}

type request struct {
	fields
}

type Balance struct {
	fields
	Fee          string `json:"fee,omitempty"`
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
	Error     Error  `json:"error,omitempty"`
}

func (a *fields) setAuthentication(key, signature string, nonce int64) {
	a.Key = key
	a.Nonce = nonce
	a.Signature = signature
}

func (a *fields) getError() error {
	e := a.Error
	if e.Message != "" {
		err := errors.New(fmt.Sprintf("%s (code: %v)", e.Message, e.Code))
		return err
	}
	return nil
}

type Error struct {
	Code    int    `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Status  int    `json:"status,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (code: %v)", e.Message, e.Code)
}

type requestBody interface {
	setAuthentication(key, signature string, nonce int64)
	getError() error
}

// Authenticate receives a Keys used to
// authenticate into the private endpoints.
func Authenticate(keys *Keys) *Account {
	return &Account{keys}
}

func (c *Account) Balance() (*Balance, error) {
	balance := &Balance{}
	if err := c.post(balancePath, balance); err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *Account) OpenOrders() ([]*Order, error) {
	var orders []*Order
	openOrders := &openOrders{}
	if err := c.post(openOrdersPath, openOrders, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (c *Account) getSignature(nonce int64) string {
	if c.validateKeys() == false {
		panic("can't generate a signature without keys")
	}
	key := c.keys.Key
	clientId := c.keys.ClientId
	secret := c.keys.Secret
	message := fmt.Sprintf("%v%v%v", nonce, key, clientId)
	signature := sign(message, secret)
	return signature
}

func (c *Account) validateKeys() bool {
	if c.keys == nil {
		return false
	}
	return true
}

func (c *Account) post(path string, schemas ...interface{}) error {
	var respSchema interface{}
	reqSchema := schemas[0].(requestBody)
	if len(schemas) == 2 {
		respSchema = schemas[1]
	} else {
		respSchema = reqSchema
	}
	nonce := getNonce()
	signature := c.getSignature(nonce)
	reqSchema.setAuthentication(c.keys.Key, signature, nonce)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, respSchema)
	if err != nil {
		respSchema = &fields{}
		if err = json.Unmarshal(body, respSchema); err != nil {
			return err
		}
	}
	if err = respSchema.(requestBody).getError(); err != nil {
		return err
	}
	return nil
}
