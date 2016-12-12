package bitso

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestClient(t *testing.T) {
	httpmock.Activate()
	registerResponder()
	defer httpmock.DeactivateAndReset()
	Convey("Given a new Client with a nil Configuration", t, func() {
		client := NewClient(nil)

		Convey("When is asked for the sandbox", func() {
			isSanbox := client.IsSandbox()

			Convey("The sanbox should be false", func() {
				So(isSanbox, ShouldBeFalse)
			})
		})

		Convey("Given the ticker path", func() {
			path := tickerPath
			ticker := &Ticker{}

			Convey("When the book is btc_mxn", func() {
				v := &url.Values{}
				v.Set("book", BTCMXN)
				err := client.get(path, v, ticker)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The price high should be 12700.00", func() {
					So(ticker.High, ShouldEqual, "12700.00")
				})
			})

			Convey("When the book is eth_mxn", func() {
				v := &url.Values{}
				v.Set("book", ETHMXN)
				err := client.get(path, v, ticker)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The price high should be 213.97", func() {
					So(ticker.High, ShouldEqual, "213.97")
				})
			})
		})

		Convey("When the ticker is requested", func() {
			Convey("And the book is btc_mxn", func() {
				ticker, err := client.Ticker(BTCMXN)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The price high should be 12700.00", func() {
					So(ticker.High, ShouldEqual, "12700.00")
				})
			})

			Convey("And the book is eth_mxn", func() {
				ticker, err := client.Ticker(ETHMXN)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The price high should be 213.97", func() {
					So(ticker.High, ShouldEqual, "213.97")
				})
			})

			Convey("And the book is invalid", func() {
				_, err := client.Ticker("invalid_book")

				Convey("An error should occur", func() {
					So(err, ShouldNotBeNil)
				})
			})
		})

		Convey("When the order book is requested", func() {
			Convey("An the book is btc_mxn", func() {
				orderBook, err := client.OrderBook(BTCMXN, false)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The bids should have length 5", func() {
					So(orderBook.Bids, ShouldHaveLength, 5)
				})
			})

			Convey("An the book is eth_mxn", func() {
				orderBook, err := client.OrderBook(ETHMXN, false)

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The bids should have length 4", func() {
					So(orderBook.Bids, ShouldHaveLength, 4)
				})
			})
		})

		Convey("When the last transactions are requested", func() {
			Convey("And the book is btc_mxn", func() {
				transactions, err := client.Transactions(BTCMXN, "")

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The transactions should have length 4", func() {
					So(transactions, ShouldHaveLength, 4)
				})

				Convey("When time is equal to minute", func() {
					transactions, err := client.Transactions(BTCMXN, "minute")

					Convey("err should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The transactions should have length 2", func() {
						So(transactions, ShouldHaveLength, 2)
					})
				})
			})

			Convey("An the book is eth_mxn", func() {
				transactions, err := client.Transactions(ETHMXN, "")

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The transactions should have length 2", func() {
					So(transactions, ShouldHaveLength, 2)
				})

				Convey("When time is equal to minute", func() {
					transactions, err := client.Transactions(ETHMXN, "minute")

					Convey("err should be nil", func() {
						So(err, ShouldBeNil)
					})

					Convey("The transactions should have length 1", func() {
						So(transactions, ShouldHaveLength, 1)
					})
				})
			})
		})
	})

	Convey("Given a Client with key, secret and id", t, func() {
		config := &Configuration{
			Key:      "key",
			Secret:   "secret",
			ClientId: "clientId",
		}
		client := NewClient(config)

		Convey("When the signature is generated", func() {
			signature := client.getSignature(getNonce())

			Convey("The signature should NOT be empty", func() {
				So(signature, ShouldNotBeEmpty)
			})

			Convey("When a new signature is generated with a new nonce", func() {
				newSignature := client.getSignature(getNonce())

				Convey("The newSignature, should be different", func() {
					So(newSignature, ShouldNotEqual, signature)
				})
			})
		})

		Convey("When the balance is requested", func() {
			balance, err := client.Balance()

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Fee should be 1.0000", func() {
				So(balance.Fee, ShouldEqual, "1.0000")
			})
		})

	})

	Convey("Given a client with an invalid key", t, func() {
		config := &Configuration{
			Key:      "",
			Secret:   "secret",
			ClientId: "clientId",
		}
		client := NewClient(config)

		Convey("When a request for the balance is made", func() {
			balance := &Balance{}
			err := client.post(balancePath, balance)

			Convey("err should be 'Invalid API Code or Invalid Signature:  (code: 101)'", func() {
				So(err.Error(), ShouldEqual, "Invalid API Code or Invalid Signature:  (code: 101)")
			})
		})

		Convey("When a request for the open orders is made", func() {
			var orders []Order
			openOrders := &OpenOrders{}
			err := client.post(openOrdersPath, openOrders, &orders)

			Convey("err should be 'Invalid API Code or Invalid Signature:  (code: 101)'", func() {
				So(err.Error(), ShouldEqual, "Invalid API Code or Invalid Signature:  (code: 101)")
			})
		})
	})

	Convey("Given an empty fields struct", t, func() {
		f := fields{}

		Convey("When the fields values are set", func() {
			f.setAuthentication("key", "signature", 1234)

			Convey("The key be equal to \"key\"", func() {
				So(f.Key, ShouldEqual, "key")
			})

			Convey("The signature be equal to \"signature\"", func() {
				So(f.Signature, ShouldEqual, "signature")
			})

			Convey("The nonce be equal to 1234", func() {
				So(f.Nonce, ShouldEqual, 1234)
			})
		})
	})
}

func registerResponder() {
	httpmock.RegisterResponder("GET", URL+tickerPath,
		func(req *http.Request) (*http.Response, error) {
			var ticker *Ticker
			v := req.URL.Query()
			book := v.Get("book")
			if book == ETHMXN {
				ticker = &Ticker{
					High:      "213.97",
					Last:      "212.30",
					Timestamp: "1468809252",
					Volume:    "149.25704647",
					Vwap:      "210.00557165",
					Low:       "205.92",
					Ask:       "212.30",
					Bid:       "208.27",
				}
			} else if book == BTCMXN || book == "" {
				ticker = &Ticker{
					High:      "12700.00",
					Last:      "12640.00",
					Timestamp: "1468809239",
					Volume:    "84.97899364",
					Vwap:      "12505.15042596",
					Low:       "12388.17",
					Ask:       "12640.00",
					Bid:       "12554.88",
				}
			}
			resp, err := httpmock.NewJsonResponse(200, ticker)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", URL+orderBookPath,
		func(req *http.Request) (*http.Response, error) {
			var orderBook *OrderBook
			v := req.URL.Query()
			book := v.Get("book")
			if book == ETHMXN {
				orderBook = &OrderBook{
					Bids: [][]string{
						[]string{
							"10720.00",
							"3.15298000",
						},
						[]string{
							"10712.40",
							"0.00326724",
						},
						[]string{
							"10711.69",
							"0.17947681",
						},
						[]string{
							"10709.96",
							"1.12340008",
						},
					},
				}
			} else if book == BTCMXN || book == "" {
				orderBook = &OrderBook{
					Bids: [][]string{
						[]string{
							"210.02",
							"2.07146938",
						},
						[]string{
							"206.62",
							"50.00000000",
						},
						[]string{
							"204.01",
							"50.00000000",
						},
						[]string{
							"204.00",
							"6.11132353",
						},
						[]string{
							"203.20",
							"10.20000000",
						},
					},
				}
			}
			resp, err := httpmock.NewJsonResponse(200, orderBook)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder("GET", URL+transactionsPath,
		func(req *http.Request) (*http.Response, error) {
			var transactions []*Transaction
			v := req.URL.Query()
			book := v.Get("book")
			time := v.Get("time")
			_ = v.Get("time")
			if book == ETHMXN {
				transactions = []*Transaction{
					&Transaction{
						Amount: "1.94511553",
						Date:   "1470876646",
						Price:  "212.03",
						Tid:    159075,
						Side:   "sell",
					},
					&Transaction{
						Amount: "1.79120536",
						Date:   "1470876493",
						Price:  "224.00",
						Tid:    159074,
						Side:   "sell",
					},
				}
				if time == "minute" {
					transactions = transactions[0:1]
				}
			} else if book == BTCMXN || book == "" {
				transactions = []*Transaction{
					&Transaction{
						Amount: "0.02200000",
						Date:   "1470876646",
						Price:  "10931.02",
						Tid:    159075,
						Side:   "sell",
					},
					&Transaction{
						Amount: "0.14089557",
						Date:   "1470876493",
						Price:  "10931.02",
						Tid:    159074,
						Side:   "sell",
					},
					&Transaction{
						Amount: "0.03561408",
						Date:   "1470876493",
						Price:  "10925.67",
						Tid:    159073,
						Side:   "sell",
					},
					&Transaction{
						Amount: "0.01737102",
						Date:   "1470876189",
						Price:  "10925.67",
						Tid:    159072,
						Side:   "sell",
					},
				}
				if time == "minute" {
					transactions = transactions[0:2]
				}
			}
			resp, err := httpmock.NewJsonResponse(200, transactions)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
	httpmock.RegisterResponder("POST", URL+balancePath,
		func(req *http.Request) (*http.Response, error) {
			balance := &Balance{
				Fee:          "1.0000",
				BTCAvailable: "46.67902107",
				MXNAvailable: "26864.57",
				BTCBalance:   "46.67902107",
				MXNBalance:   "26864.57",
				BTCReserved:  "0.00000000",
				MXNReserved:  "0.00",
			}
			r := req.Body
			body, err := ioutil.ReadAll(r)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			err = json.Unmarshal(body, balance)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			if balance.Key != "key" {
				e := Error{
					Code:    101,
					Message: "Invalid API Code or Invalid Signature: " + balance.Key,
				}
				balance.Error = e
			}
			resp, err := httpmock.NewJsonResponse(200, balance)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	httpmock.RegisterResponder("POST", URL+openOrdersPath,
		func(req *http.Request) (*http.Response, error) {
			openOrders := &OpenOrders{}
			orders := []*Order{
				&Order{
					Amount:   "0.01000000",
					Datetime: "2015-11-12 12:37:01",
					Price:    "5600.00",
					Id:       "543cr2v32a1h684430tvcqx1b0vkr93wd694957cg8umhyrlzkgbaedmf976ia3v",
					Type:     "1",
					Status:   "1",
				},
				&Order{
					Amount:   "0.12680000",
					Datetime: "2015-11-12 12:33:47",
					Price:    "4000.00",
					Id:       "qlbga6b600n3xta7actori10z19acfb20njbtuhtu5xry7z8jswbaycazlkc0wf1",
					Type:     "0",
					Status:   "0",
				},
				&Order{
					Amount:   "1.12560000",
					Datetime: "2015-11-12 12:33:23",
					Price:    "6123.55",
					Id:       "d71e3xy2lowndkfmde6bwkdsvw62my6058e95cbr08eesu0687i5swyot4rf2yf8",
					Type:     "1",
					Status:   "0",
				},
			}
			r := req.Body
			body, err := ioutil.ReadAll(r)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			err = json.Unmarshal(body, openOrders)
			if err != nil {
				return httpmock.NewStringResponse(500, err.Error()), nil
			}
			if openOrders.Key != "key" {
				e := Error{
					Code:    101,
					Message: "Invalid API Code or Invalid Signature: " + openOrders.Key,
				}
				f := fields{
					Error: e,
				}
				return httpmock.NewJsonResponse(200, f)
			}
			resp, err := httpmock.NewJsonResponse(200, orders)
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)
}
