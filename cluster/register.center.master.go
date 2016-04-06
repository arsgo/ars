package cluster

func (d *rcServer) Bind() (err error) {
	err = d.createRCServer(d.rcServerPath,d.dataMap.Translate(rcServerValue))
	if err != nil {
		d.Log.Error(err)
		return
	}
	go d.watchMasterChange()
	go d.resetRCSnap()
	return
}
