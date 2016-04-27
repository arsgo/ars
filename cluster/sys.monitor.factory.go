package cluster

import (
	"errors"
	"fmt"

	"github.com/colinyl/ars/mqservice"
)

type monitorHandler interface {
	Send(string, string) error
	Close()
}

func getMonitorHandler(typeName string, content string) (monitorHandler, error) {
	switch typeName {
	case "mq":
		return mqservice.NewMQService(content), nil
	}
	return nil, errors.New(fmt.Sprintf("not support: %s", typeName))
}
