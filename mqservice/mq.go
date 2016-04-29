package mqservice

import "github.com/colinyl/stomp"

type IMQService interface {
	Consume(string, func(stomp.MsgHandler)bool) error
	Send(string, string) error
	Close()
}
