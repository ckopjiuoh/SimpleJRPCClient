package main

import (
	"APITestingFramework/jrpc-client"
	log "github.com/sirupsen/logrus"
)

type CalcAdd struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func init(){
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {


	client := jrpc_client.NewClient("https://gurujsonrpc.appspot.com/guru")
	client.Call("guru.test", [1]string{"Guru"})
	client.Call("guru.f", [1]string{"Guru"})
	client.Call("guru.d", [1]string{"Guru"})
	client.Call("guru.add", [1]string{"Guru"})
	client.Call("guru.de", [1]string{"Guru"})
}
