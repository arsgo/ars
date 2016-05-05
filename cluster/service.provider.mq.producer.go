package cluster

func (d *spServer) NewMQProducer(name string) (p *MQProducer, err error) {
	config, err := d.zkClient.GetMQConfig(name)
	if err != nil {
		return
	}
	p, err = NewMQProducer(config)
	return
}
