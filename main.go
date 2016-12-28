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
	ticker, err := bitso.Ticker(bitso.BTCMXN)
	if err != nil {
		panic(err)
	}
	fmt.Println(ticker)
	orderBook, err := bitso.OrderBook(bitso.BTCMXN, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(orderBook)
	transactions, err := bitso.Transactions(bitso.BTCMXN, "hour")
	if err != nil {
		panic(err)
	}
	fmt.Println(transactions)
	balance, err := account.Balance()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("balance", balance)
	orders, err := account.OpenOrders()
	if err != nil {
		panic(err)
	}
	order := orders[0]
	orders, err = account.LookupOrder(order.Id)
	if err != nil {
		panic(err)
	}
	fmt.Println("orders", orders)
}
