package jrpc_client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"errors"

	"SimpleJRPCClient/model"
	"io/ioutil"

	"log"
	"fmt"
)

type JRPCClient struct {
	Addr    string
	Method  string
	Headers map[string]string
	Body    model.RPCRequestBody
	Count int

	httpClient *http.Client
}

func NewClient(addr string) *JRPCClient {

	return &JRPCClient{
		Addr: addr,
		httpClient: &http.Client{
		},
		Method:  "Post",
		Headers: map[string]string{"Content-type": "application/json"},
		Count: 0,
	}
}

func (c *JRPCClient) WithHeaders(m map[string]string) *JRPCClient {
	c.Headers = m

	return c
}

func (c *JRPCClient) WithMethod(m string) *JRPCClient {
	c.Method = m

	return c
}

func (c *JRPCClient) WithBody(method string, params interface{}) *JRPCClient {
	var m = model.RPCRequestBody{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	c.Body = m
	return c
}

func (c *JRPCClient) Call() (*model.RPCResponse, error) {

	client := c.httpClient
	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(&c.Body)

	req, err := http.NewRequest(c.Method, c.Addr, body)

	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	if err != nil {
		log.Fatalln(err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer resp.Body.Close()
	r := new(model.RPCResponse)
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&r)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Respnse decoding error: %s", err.Error()))
		if e, ok := err.(*json.SyntaxError); ok {
			log.Fatalln(fmt.Sprintf("Syntax error: %s, \n Offset error: %s", e.Error(), e.Offset))
		}

		body, _ := ioutil.ReadAll(resp.Body)

		log.Fatalln("RPC response: %q", string(body))
		return nil, errors.New("Not a RPC Response!")
	}
	c.Count++
	return r, nil
}
