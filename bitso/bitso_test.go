package bitso

import (
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBitso(t *testing.T) {
	httpmock.Activate()
	registerResponder()
	defer httpmock.DeactivateAndReset()
	Convey("When the ticker is requested", t, func() {
		Convey("And the book is btc_mxn", func() {
			ticker, err := Ticker(BTCMXN)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The price high should be 12700.00", func() {
				So(ticker.High, ShouldEqual, "12700.00")
			})
		})

		Convey("And the book is eth_mxn", func() {
			ticker, err := Ticker(ETHMXN)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The price high should be 213.97", func() {
				So(ticker.High, ShouldEqual, "213.97")
			})
		})

		Convey("And the book is invalid", func() {
			_, err := Ticker("invalid_book")

			Convey("An error should occur", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("When the order book is requested", t, func() {
		Convey("An the book is btc_mxn", func() {
			orderBook, err := OrderBook(BTCMXN, false)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The bids should have length 5", func() {
				So(orderBook.Bids, ShouldHaveLength, 5)
			})
		})

		Convey("An the book is eth_mxn", func() {
			orderBook, err := OrderBook(ETHMXN, false)

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The bids should have length 4", func() {
				So(orderBook.Bids, ShouldHaveLength, 4)
			})
		})
	})

	Convey("When the last transactions are requested", t, func() {
		Convey("And the book is btc_mxn", func() {
			transactions, err := Transactions(BTCMXN, "")

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The transactions should have length 4", func() {
				So(transactions, ShouldHaveLength, 4)
			})

			Convey("When time is equal to minute", func() {
				transactions, err := Transactions(BTCMXN, "minute")

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The transactions should have length 2", func() {
					So(transactions, ShouldHaveLength, 2)
				})
			})
		})

		Convey("An the book is eth_mxn", func() {
			transactions, err := Transactions(ETHMXN, "")

			Convey("err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The transactions should have length 2", func() {
				So(transactions, ShouldHaveLength, 2)
			})

			Convey("When time is equal to minute", func() {
				transactions, err := Transactions(ETHMXN, "minute")

				Convey("err should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The transactions should have length 1", func() {
					So(transactions, ShouldHaveLength, 1)
				})
			})
		})
	})

	Convey("Given a unique nonce", t, func() {
		nonce := getNonce()

		Convey("When a new nonce is generated", func() {
			newNonce := getNonce()

			Convey("The new nonce should be greater than the previous one", func() {
				So(newNonce, ShouldBeGreaterThan, nonce)
			})
		})
	})

	Convey("Given a message to sign with a private key", t, func() {
		message := "message"
		key := "secret"

		Convey("When the message signed", func() {
			signature := sign(message, key)

			Convey("The signature should be as expected", func() {
				So(signature, ShouldEqual, "8b5f48702995c1598c573db1e21866a9b825d4a794d169d7060a03605796360b")
			})
		})
	})
}

func Test_validateBook(t *testing.T) {
	Convey("Given a book to validate", t, func() {
		var validation bool

		Convey("When the book is btc_mxn", func() {
			validation = validateBook(BTCMXN)

			Convey("validateBook should be true", func() {
				So(validation, ShouldBeTrue)
			})
		})

		Convey("When the book is eth_mxn", func() {
			validation = validateBook(ETHMXN)

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
