package bitso

import (
  "testing"
  . "github.com/smartystreets/goconvey/convey"
)

func TestBitso(t *testing.T) {
  Convey("Given a url to GET", t, func() {
    path := tickerPath

    Convey("When the request is done", func() {
      ticker := &Ticker{}
      body, err := get(path, nil, ticker)
      So(err, ShouldBeNil)
      So(body, ShouldNotBeNil)
    })
  })
}
