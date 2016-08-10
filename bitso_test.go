package bitso

import (
  "net/url"
  "net/http"
  "testing"
  . "github.com/smartystreets/goconvey/convey"
  "github.com/jarcoal/httpmock"
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
        v.Set("book", btcmxn)
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
        v.Set("book", ethmxn)
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
        ticker, err := client.Ticker(btcmxn)

        Convey("err should be nil", func() {
          So(err, ShouldBeNil)
        })

        Convey("The price high should be 12700.00", func() {
          So(ticker.High, ShouldEqual, "12700.00")
        })
      })

      Convey("And the book is eth_mxn", func() {
        ticker, err := client.Ticker(ethmxn)

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
        orderBook, err := client.OrderBook(btcmxn, false)

        Convey("err should be nil", func() {
          So(err, ShouldBeNil)
        })

         Convey("The bids should have length 5", func() {
           So(orderBook.Bids, ShouldHaveLength, 5)
         })
      })

      Convey("An the book is eth_mxn", func() {
        orderBook, err := client.OrderBook(ethmxn, false)

        Convey("err should be nil", func() {
          So(err, ShouldBeNil)
        })

        Convey("The bids should have length 4", func() {
          So(orderBook.Bids, ShouldHaveLength, 4)
        })
      })
    })
  })
}

func Test_validateBook(t *testing.T) {
  Convey("Given a book to validate", t, func() {
    var validation bool

    Convey("When the book is btc_mxn", func() {
      validation = validateBook(btcmxn)

      Convey("validateBook should be true", func() {
        So(validation, ShouldBeTrue)
      })
    })

    Convey("When the book is eth_mxn", func() {
      validation = validateBook(ethmxn)

      Convey("validateBook should be true", func() {
        So(validation, ShouldBeTrue)
      })
    })

    Convey("When the book is invalid", func() {
      validation = validateBook("invalid_book")

      Convey("validateBook should be false", func() {
        So(validation, ShouldBeFalse)
      })
    })
  })
}

func registerResponder() {
  httpmock.RegisterResponder("GET", URL + tickerPath,
    func(req *http.Request) (*http.Response, error) {
      var ticker *Ticker
      v := req.URL.Query()
      book := v.Get("book")
      if book == ethmxn {
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
      } else if book == btcmxn || book == "" {
        ticker = &Ticker{
          High: "12700.00",
          Last: "12640.00",
          Timestamp: "1468809239",
          Volume: "84.97899364",
          Vwap: "12505.15042596",
          Low: "12388.17",
          Ask: "12640.00",
          Bid: "12554.88",
        }
      }
      resp, err := httpmock.NewJsonResponse(200, ticker)
      if err != nil {
        return httpmock.NewStringResponse(500, ""), nil
      }
      return resp, nil
    },
  )

  httpmock.RegisterResponder("GET", URL + orderBookPath,
    func(req *http.Request) (*http.Response, error) {
      var orderBook *OrderBook
      v := req.URL.Query()
      book := v.Get("book")
      if book == ethmxn {
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
      } else if book == btcmxn || book == "" {
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
}
