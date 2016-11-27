package main

import (
	"fmt"
	"os"

	"github.com/dsmontoya/gobitso/bitso"
)

func main() {
	//CLI
	conf := &bitso.Configuration{
		Key:      os.Getenv("BITSO_KEY"),
		Secret:   os.Getenv("BITSO_SECRET"),
		ClientId: os.Getenv("BITSO_CLIENT_ID"),
	}
	client := bitso.NewClient(conf)
	ticker, err := client.Ticker(bitso.BTCMXN)
	if err != nil {
		panic(err)
	}
	fmt.Println(ticker)
	orderBook, err := client.OrderBook(bitso.BTCMXN, false)
	if err != nil {
		panic(err)
	}
	fmt.Println(orderBook)
	transactions, err := client.Transactions(bitso.BTCMXN, "hour")
	if err != nil {
		panic(err)
	}
	fmt.Println(transactions)
	balance, err := client.Balance()
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("balance", balance)
	openOrders, err := client.OpenOrders()
	if err != nil {
		panic(err)
	}
	fmt.Println(openOrders)
}
