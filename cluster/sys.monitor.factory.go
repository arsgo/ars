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

func StaticSendMonitor(typeName string, config string, queue string, content string) (err error) {
	handler, err := getMonitorHandler(typeName, config)
	if err != nil {
		return
	}
	err = handler.Send(queue, content)
	handler.Close()
	return

}

func getMonitorHandler(typeName string, content string) (monitorHandler, error) {
	switch typeName {
	case "mq":
		return mqservice.NewMQService(content), nil
	}
	return nil, errors.New(fmt.Sprintf("not support: %s", typeName))
}
