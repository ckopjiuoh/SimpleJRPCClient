package test

import (
	. "SimpleJRPCClient/jrpc-client"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestClient(t *testing.T) {
	var client = NewClient("https://gurujsonrpc.appspot.com/guru")
	Convey("Should be POST guru.test method", t, func() {
		r, _ := client.
			WithMethod("POST").
			WithBody("guru.test", []string{"guru"}).
			Call()
		Convey("When response result is 2", func() {
			So(r.Result.(string), ShouldContainSubstring, "Hello guru!")
		})
		Convey("When error is nil", func() {
			So(r.Error, ShouldBeNil)
		})
	})
}
