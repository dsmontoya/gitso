package bitso

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	URL              = "https://api.bitso.com/v2/"
	BTCMXN           = "btc_mxn"
	ETHMXN           = "eth_mxn"
	tickerPath       = "ticker"
	transactionsPath = "transactions"
	orderBookPath    = "order_book"
	balancePath      = "balance"
	openOrdersPath   = "open_orders"
)

type TickerInfo struct {
	High      string
	Last      string
	Timestamp string
	Volume    string
	Vwap      string
	Low       string
	Ask       string
	Bid       string
}

type OrderBookInfo struct {
	Asks [][]string
	Bids [][]string
}

type Transaction struct {
	Amount string
	Date   string
	Price  string
	Tid    int
	Side   string
}

func Ticker(book string) (*TickerInfo, error) {
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	ticker := &TickerInfo{}
	v := &url.Values{}
	v.Set("book", book)
	err := get(tickerPath, v, ticker)
	if err != nil {
		return nil, err
	}
	return ticker, nil
}

func OrderBook(book string, group bool) (*OrderBookInfo, error) {
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	orderBook := &OrderBookInfo{}
	v := &url.Values{}
	v.Set("book", book)
	err := get(orderBookPath, v, orderBook)
	if err != nil {
		return nil, err
	}
	return orderBook, nil
}

/*
GetTransactions returns a list of recent trades from the specified book
and the specified time frame.

Valid time frames are hour and minute. Leaving time blank will set hour as the default frame.
*/
func Transactions(book string, time string) ([]*Transaction, error) {
	var transactions []*Transaction
	if validateBook(book) == false {
		err := errors.New("Invalid book value")
		return nil, err
	}
	v := &url.Values{}
	v.Set("book", book)
	v.Set("time", time)
	err := get(transactionsPath, v, &transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func getNonce() int64 {
	return time.Now().UnixNano()
}

func validateBook(book string) bool {
	return book == BTCMXN || book == ETHMXN
}

func sign(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	bytes := mac.Sum(nil)
	s := hex.EncodeToString(bytes)
	return s
}

func get(path string, query *url.Values, schema interface{}) error {
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
