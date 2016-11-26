package bitso

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	URL              = "https://api.bitso.com/v2/"
	btcmxn           = "btc_mxn"
	ethmxn           = "eth_mxn"
	tickerPath       = "ticker"
	transactionsPath = "transactions"
	orderBookPath    = "order_book"
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
	return book == btcmxn || book == ethmxn
}

func sign(message, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(message))
	bytes := mac.Sum(nil)
	s := hex.EncodeToString(bytes)
	return s
}
