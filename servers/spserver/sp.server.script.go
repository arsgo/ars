package main

import "github.com/arsgo/lib4go/db"

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
func (s *SPServer) CloseDB() {
	dbs := s.dbPool.GetAll()
	for _, d := range dbs {
		cdb := d.(*db.DBScriptBind)
		cdb.Close()
	}
}
