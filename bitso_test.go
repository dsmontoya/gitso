package bitso

import (
  "net/url"
  "net/http"
  "testing"
  . "github.com/smartystreets/goconvey/convey"
  "github.com/jarcoal/httpmock"
)

func TestRequest(t *testing.T) {
  httpmock.Activate()
  registerResponder()
  defer httpmock.DeactivateAndReset()
  Convey("Given the ticker path", t, func() {
    path := tickerPath
    ticker := &Ticker{}

    Convey("When the book is btc_mxn", func() {
      v := &url.Values{}
      v.Set("book", btcmxn)
      err := get(path, v, ticker)

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
      err := get(path, v, ticker)

      Convey("err should be nil", func() {
        So(err, ShouldBeNil)
      })

      Convey("The price high should be 213.97", func() {
        So(ticker.High, ShouldEqual, "213.97")
      })
    })
  })
}

func TestBitso(t *testing.T) {
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
  })
  Convey("Public methods", t, func() {

  })
}

func registerResponder() {
  httpmock.RegisterResponder("GET", API_URL + tickerPath,
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
}
