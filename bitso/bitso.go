package bitso

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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

type Ticker struct {
	High      string
	Last      string
	Timestamp string
	Volume    string
	Vwap      string
	Low       string
	Ask       string
	Bid       string
}

type OrderBook struct {
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
