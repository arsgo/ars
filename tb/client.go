package main

import (
	"strings"
	"time"

	"github.com/colinyl/ars/rpcservice"
)

type TCPClient struct {
	client  *rpcservice.RPCClient
	address string
}


func NewTCPClient(address string) *TCPClient {
	client := &TCPClient{client: rpcservice.NewRPCClient(address), address: address}
	client.client.Open()
	return client
}

func (c *TCPClient) Reqeust() (resp *response) {
	defer func() {
		if err := recover(); nil != err {
			Log.Fatal(err.(error).Error())
			resp = &response{success: false, url: c.address, useTime: 0}
		}
	}()
	startTime := time.Now()

	result, err := c.client.Request("test_request", "{}")
	if err != nil {
		Log.Print(err)
	}
	// Log.Info(result)
	endTime := time.Now()
	return &response{success: err == nil && strings.EqualFold(result, "success"), url: c.address, useTime: subTime(startTime, endTime)}

}
func subTime(startTime time.Time, endTime time.Time) int {
	return int(endTime.Sub(startTime).Nanoseconds() / 1000 / 1000)
}
func (c *TCPClient) Close(){
    c.client.Close()
}