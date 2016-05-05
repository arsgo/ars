package cluster

import "strings"

type serviceProviderConsumer struct {
	sp *spServer
}

func (c *serviceProviderConsumer) GetSourceConfig(typeName string, name string) (config string, err error) {
	return c.sp.zkClient.GetSourceConfig(typeName, name)
}
func NewServiceProviderConsumer(sp *spServer) *serviceProviderConsumer {
	return &serviceProviderConsumer{sp: sp}
}
func (d *spServer) setMQConsumer(services map[string]spService) {
	data := make(map[string]spService)
	for i, v := range services {
		if strings.EqualFold(strings.ToLower(v.Type), "mq") && strings.EqualFold(strings.ToLower(v.Method), "consumer") {
			data[i] = v
		}
	}
	d.mqConsumerManager.Reset(data)
}

func (c *serviceProviderConsumer) Hande(script string, params string, h *msgHandler) bool {
	values, err := c.sp.scriptEngine.script.Call(script, getScriptInputArgs(h.Message, params))
	return err != nil && strings.EqualFold(values[0], "true")
}
