package cluster

import (
	"github.com/colinyl/lib4go/influxdb"
	"github.com/colinyl/lib4go/mem"
)

func (d *spServer) NewInfluxDB(name string) (p *influxdb.InfluxDB, err error) {
	config, err := d.zkClient.GetSourceConfig("db", name)
	if err != nil {
		return
	}
	p, err = influxdb.New(config)
	return
}

func (d *spServer) NewMemcached(name string) (p *mem.MemcacheClient, err error) {
	config, err := d.zkClient.GetSourceConfig("db", name)
	if err != nil {
		return
	}
	p, err = mem.New(config)
	return
}
