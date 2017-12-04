package test

import (
	. "SimpleJRPCClient/jrpc-client"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestClient(t *testing.T) {
	var client = NewClient("https://gurujsonrpc.appspot.com/guru")
	Convey("Given POST guru.test with params [guru]", t, func() {
		r, _ := client.
			WithMethod("POST").
			WithBody("guru.test", []string{"guru"}).
			Call()
		Convey("The response result should be 2", func() {
			So(r.Result.(string), ShouldContainSubstring, "Hello guru!")
		})
		Convey("The error should be nil", func() {
			So(r.Error, ShouldBeNil)
		})
	})
}
