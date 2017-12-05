package test

import (
	. "SimpleJRPCClient/jrpc-client"
	"testing"
	"encoding/json"
)

func TestClient(t *testing.T) {
	var client = NewClient("https://gurujsonrpc.appspot.com/guru")
		r, _ := client.
			WithMethod("POST").
			WithBody("guru.test", []string{"guru"}).
			Call()
			var s string
			json.Unmarshal(*r.Result, &s)
}
