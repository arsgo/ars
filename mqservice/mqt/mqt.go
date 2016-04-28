package main

import (
	"fmt"
	"time"

	"github.com/colinyl/ars/mqservice"
	zk "github.com/colinyl/lib4go/zkclient"
)

func main() {
	zk.New([]string{"192.168.101.161:2181"}, time.Second)

	for i := 1; i < 10000; i++ {
		fmt.Println(i)
		stomp := mqservice.NewStompService(`{"type": "stomp","address": "192.168.101.161:61613"}`)
		stomp.Send("go:t:qu", string(time.Now().Unix()))
		stomp.Close()
		time.Sleep(time.Second)
	}

}
