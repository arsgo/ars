package cluster

import (
	"github.com/colinyl/lib4go/elastic"
	"github.com/colinyl/lib4go/influxdb"
	"github.com/colinyl/lib4go/mem"
)

func (d *clusterClient) NewInfluxDB(name string) (p *influxdb.InfluxDB, err error) {
	config, err := d.GetSourceConfig("db", name)
	if err != nil {
		return
	}
	p, err = influxdb.New(config)
	return
}

func (d *clusterClient) NewMemcached(name string) (p *mem.MemcacheClient, err error) {
	config, err := d.GetSourceConfig("db", name)
	if err != nil {
		return
	}
	p, err = mem.New(config)
	return
}

func (d *clusterClient) NewMQProducer(name string) (p *MQProducer, err error) {
	config, err := d.GetMQConfig(name)
	if err != nil {
		return
	}
	p, err = NewMQProducer(config)
	return
}

func (d *clusterClient) NewElastic(name string) (es *elastic.ElasticSearch, err error) {
	config, err := d.GetElasticConfig(name)
	if err != nil {
		return
	}
	es, err = elastic.New(config)
	return
}
