package bitso

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBitso(t *testing.T) {
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
