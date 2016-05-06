package cluster

import (
	"encoding/json"

	"github.com/colinyl/lib4go/elastic"
)

type elasticConfig struct {
	Host []string `json:"hosts"`
}

func (d *spServer) NewElastic(name string) (es *elastic.ElasticSearch, err error) {
	config, err := d.zkClient.GetElasticConfig(name)
	if err != nil {
		return
	}
	var elasticonfig elasticConfig
	err = json.Unmarshal([]byte(config), &elasticonfig)
	if err != nil {
		return
	}
	es = elastic.New(elasticonfig.Host)
	return
}
