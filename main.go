package main

import (
	"fmt"
	"os"

	"github.com/dsmontoya/gobitso/bitso"
)

func main() {
	//CLI
	keys := &bitso.Keys{
		Key:      os.Getenv("BITSO_KEY"),
		Secret:   os.Getenv("BITSO_SECRET"),
		ClientId: os.Getenv("BITSO_CLIENT_ID"),
	}
	account := bitso.Authenticate(keys)
	ticker, err := bitso.GetTicker(bitso.BTCMXN)
	if err != nil {
		panic(err)
	}
	fmt.Println(ticker)
	orderBook, err := bitso.GetOrderBook(bitso.BTCMXN, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(orderBook)
	transactions, err := bitso.GetTransactions(bitso.BTCMXN, "hour")
	if err != nil {
		panic(err)
	}
	fmt.Println(transactions)
	balance, err := account.GetBalance()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("balance", balance)
	openOrders, err := account.GetOpenOrders()
	if err != nil {
		panic(err)
	}
	fmt.Println(openOrders)
}
