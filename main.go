package main

import (
	"APITestingFramework/jrpc-client"
	"github.com/sirupsen/logrus"
	"os"
)

func init(){
	jrpc_client.Log = logrus.New()
	jrpc_client.Log.Formatter = &logrus.JSONFormatter{}
	jrpc_client.Log.Out = os.Stdout
	file, _ := os.OpenFile("logrus.Log", os.O_CREATE|os.O_WRONLY, 0666)
	jrpc_client.Log.Out = file

}

func main() {
	client := jrpc_client.NewClient("https://gurujsonrpc.appspot.com/guru")
	client.Call("guru.test", [1]string{"Gur"})
	client.Call("guru.test", [1]string{"Guu"})
	client.Call("guru.test", [1]string{"Gru"})
	client.Call("guru.test", [1]string{"uru"})
	client.Call("guru.test", [1]string{"Gurui"})
}

