package mqservice

import "github.com/colinyl/stomp"

type IMQService interface {
	Consume(string, func(stomp.MsgHandler)) error
	Send(string, string) error
	Close()
}
