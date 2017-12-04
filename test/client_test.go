package test

import (
	. "SimpleJRPCClient/jrpc-client"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var client = NewClient(
	"https",
	"gurujsonrpc.appspot.com",
	0,
	"guru").
	WithPost().
	WithRPCMethod("guru.test")

func TestClient(t *testing.T) {
	Convey("Should be POST guru.test method", t, func() {
		r, _ := client.
			WithRPCParams([]string{"guru"}).
			Call()
		Convey("The response result should be 2", func() {
			var s string
			json.Unmarshal(*r.Result, &s)
			So(s, ShouldContainSubstring, "Hello guru!")
		})
		Convey("The error should be nil", func() {
			So(r.Error, ShouldBeNil)
		})
	})
}
