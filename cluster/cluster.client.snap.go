package cluster

func (d *ClusterClient) ResetSnap(addr string, snap string) (err error) {
	if d.handler.Exists(addr) {
		err = d.handler.UpdateValue(addr, snap)
	} else {
		_, err = d.handler.CreateTmpNode(addr, snap)
	}
	return
}
func (d *ClusterClient) ResetAppServerSnap(snap string) (err error) {
	return d.ResetSnap(d.appServerPath, snap)
}
