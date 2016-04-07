package cluster

//WatchMasterChange  watch whether the master server is changed
func (d *rcServer) watchMasterChange() {
	if d.IsMasterServer {
		go d.watchServiceProviderChange()
		return
	}
   d.Log.Info("::watch for master server changes")
	watchZKChildrenPathChange(d.rcServerRoot, func() {
		if m := d.isMaster(); m && !d.IsMasterServer {
			d.setOnlineParams(true)
			d.resetRCSnap()
			d.watchMasterChange()
		}
	})
}
