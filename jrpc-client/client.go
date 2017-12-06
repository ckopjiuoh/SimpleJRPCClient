package jrpc_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ckopjiuoh/SimpleJRPCClient/model"
	"io/ioutil"

	"fmt"
	"log"
)

type JRPCClient struct {
	Addr    string
	Method  string
	Headers map[string]string
	Body    model.RPCRequestBody
	Count   int

	httpClient *http.Client
}

func NewClient(protocol string, host string, port int, uri string) *JRPCClient {

	var addr = fmt.Sprintf("%s://%s", protocol, host)
	if port > 0 {
		addr = fmt.Sprintf("%s:%d", addr, port)
	}
	if uri != "" {
		addr = fmt.Sprintf("%s/%s", addr, uri)
	}

	return &JRPCClient{
		Addr:       addr,
		httpClient: &http.Client{},
		Method:     "POST",
		Headers:    map[string]string{"Content-type": "application/json"},
		Body: model.RPCRequestBody{
			Jsonrpc: "2.0",
			ID:      1,
		},
		Count: 0,
	}
}

func (c *JRPCClient) WithHeaders(m map[string]string) *JRPCClient {
	c.Headers = m

	return c
}

func (c *JRPCClient) WithPost() *JRPCClient {
	c.Method = "POST"

	return c
}

func (c *JRPCClient) WithRPCMethod(method string) *JRPCClient {
	c.Body.Method = method
	return c
}

func (c *JRPCClient) WithRPCParams(params interface{}) *JRPCClient {
	c.Body.Params = params
	return c
}

func (c *JRPCClient) Call() (*model.RPCResponse, error) {
	data, err := c.Body.MarshalJSON()
	client := c.httpClient
	req, err := http.NewRequest(c.Method, c.Addr, bytes.NewReader(data))

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
	r := model.RPCResponse{}
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = r.UnmarshalJSON(data)
	if err == nil {
		return &r, nil
	}

	log.Fatalln(fmt.Sprintf("Response unmarshalling error: %s", err))
	if e, ok := err.(*json.SyntaxError); ok {
		log.Fatalln(fmt.Sprintf("Syntax error: %s, \n Offset error: %s", e.Error(), e.Offset))
	}

	body, _ := ioutil.ReadAll(resp.Body)

	log.Fatalln("RPC response: %q", string(body))
	return nil, errors.New("Not a RPC Response!")

}
