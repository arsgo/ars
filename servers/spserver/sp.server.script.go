package main

import (
	"github.com/arsgo/lib4go/db"
	"github.com/arsgo/lib4go/elastic"
	"github.com/arsgo/lib4go/influxdb"
)

//GetScriptBinder 获取脚本绑定
func (s *SPServer) GetScriptBinder() (funs map[string]interface{}) {
	funs = map[string]interface{}{
		"NewInfluxDB": s.NewInfluxDB,
		"NewElastic":  s.NewElastic,
		"NewDB":       s.NewDB,
	}
	return
}

//NewInfluxDB 创建InfluxDB操作对象
func (s *SPServer) NewInfluxDB(name string) (p *influxdb.InfluxDB, err error) {
	config, err := s.clusterClient.GetDBConfig(name)
	if err != nil {
		return
	}
	p, err = influxdb.New(config)
	return
}

//NewElastic 创建Elastic对象
func (s *SPServer) NewElastic(name string) (es *elastic.ElasticSearch, err error) {
	config, err := s.clusterClient.GetElasticConfig(name)
	if err != nil {
		return
	}
	es, err = elastic.New(config)
	return
}

//NewDB NewDB
func (s *SPServer) NewDB(name string) (bind *db.DBScriptBind, err error) {
	config, err := s.clusterClient.GetDBConfig(name)
	if err != nil {
		return
	}
	p := s.dbPool.Get(name)
	if p != nil {
		bind = p.(*db.DBScriptBind)
		bind.ResetPoolSize(config)
		return
	}
	bind, err = db.NewDBScriptBind(config)
	if err != nil {
		return
	}
	s.dbPool.Set(name, bind)
	return
}
