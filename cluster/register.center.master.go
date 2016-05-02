package cluster

func (d *rcServer) Bind() (err error) {
	err = d.createRCServer(d.rcServerPath,d.snap.GetSnap())
	if err != nil {
		d.Log.Error(err)
		return
	}
	go d.watchMasterChange()
	go d.ResetRCSnap()
	return
}
