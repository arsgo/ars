package main

import (
	"encoding/json"
	"runtime/debug"
	"strings"
	"time"

	"github.com/arsgo/ars/rpcservice"
	"github.com/arsgo/lib4go/utility"
)

type resultCode struct {
	Code string `json:"code"`
}

type TCPClient struct {
	client      *rpcservice.RPCClient
	address     string
	params      string
	commandName string
}

func NewTCPClient(address string, commandName string, params string) *TCPClient {
	client := &TCPClient{address: address, params: params, commandName: commandName}
	client.client = rpcservice.NewRPCClientTimeout(address, time.Second*5, "tb")
	err := client.client.Open()
	if err != nil {
		Log.Fatal(err)
	}
	return client
}

func (c *TCPClient) Reqeust() (resp *response) {
	defer func() {
		if err := recover(); nil != err {
			Log.Fatal(err, string(debug.Stack()))
			resp = &response{success: false, url: c.address, useTime: 0}
		}
	}()

	startTime := time.Now()

	/*if err != nil {
		Log.Println(err)
		resp = &response{success: false, url: c.address, useTime: 0}
		return
	}*/
	c.client.Open()
	result, err := c.client.Request(c.commandName, c.params, utility.GetSessionID())
	defer c.client.Close()
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
		Log.Print(err, "[", result, "]")
	} else if !strings.EqualFold(code.Code, "success") {
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
