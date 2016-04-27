package mqt

import (
	"fmt"
	"time"

	"github.com/colinyl/ars/mqservice"
)

func main() {
	stomp := mqservice.NewStompService(`{"type": "stomp","address": "192.168.101.161:61613"}`)
	for i := 1; i < 10000; i++ {
		fmt.Println(i)
		stomp.Send("go:t:qu", string(time.Now().Unix()))
		time.Sleep(time.Second)
	}

}
