package main

import (
	"fmt"
	"os"
)

type TCPClients struct {
	clients []*TCPClient
	address string
	count   int
}

func NewHTCPClients(count int, address string) *TCPClients {
	c := &TCPClients{count: count,address:address}
	c.clients = make([]*TCPClient, 0)
	for i := 0; i < c.count; i++ {
		c.clients = append(c.clients, NewTCPClient(address))
	}
	return c
}

func (c *TCPClients) RunNow(i int) *response {
	if i > len(c.clients)-1 {
		fmt.Printf("索引错误:%d\r\n", i)
		os.Exit(1)
	}
	client := c.clients[i]
	return client.Reqeust()
}
func (c *TCPClients) GetLen() int {
	return len(c.clients)
}
func (c *TCPClients) Close(){
    for _,v:=range c.clients{
        v.Close()
    }
}