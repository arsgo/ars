package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/colinyl/ars/rpcservice"
)

type resultCode struct {
	Code string `json:"code"`
}

type TCPClient struct {
	client  *rpcservice.RPCClient
	address string
}

func NewTCPClient(address string) *TCPClient {
	client := &TCPClient{address: address}
	client.client = rpcservice.NewRPCClientTimeout(address, 10)
	err := client.client.Open()
	if err != nil {
		Log.Fatal(err)
	}
	return client
}

func (c *TCPClient) Reqeust() (resp *response) {
	defer func() {
		//if c.client != nil {
		//c.client.Close()
		//}
		if err := recover(); nil != err {
			Log.Fatal(err.(error).Error())
			resp = &response{success: false, url: c.address, useTime: 0}
		}
	}()

	startTime := time.Now()

	/*if err != nil {
		Log.Println(err)
		resp = &response{success: false, url: c.address, useTime: 0}
		return
	}*/
	result, err := c.client.Request("save_logger", "{}")
	if err != nil {
		Log.Print(err)
	}
	time.Sleep(time.Millisecond / 10)

	// Log.Info(result)
	endTime := time.Now()
	code := &resultCode{}
	err = json.Unmarshal([]byte(result), &code)
	var isSuccess bool
	if err != nil {
		Log.Print(err)
	} else if !strings.EqualFold(code.Code, "100") {
		Log.Print(result)
	} else {
		isSuccess = true
	}

	return &response{success: err == nil && isSuccess, url: c.address, useTime: subTime(startTime, endTime)}

}
func subTime(startTime time.Time, endTime time.Time) int {
	return int(endTime.Sub(startTime).Nanoseconds() / 1000 / 1000)
}
func (c *TCPClient) Close() {
	c.client.Close()
}
