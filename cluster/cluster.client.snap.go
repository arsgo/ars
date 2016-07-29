package cluster

func (d *ClusterClient) UpdateSnap(addr string, snap string) (err error) {
	if d.handler.Exists(addr) {
		err = d.handler.UpdateValue(addr, snap)
	} else {
		_, err = d.handler.CreateTmpNode(addr, snap)
	}
	return
}
func (d *ClusterClient) UpdateAppServerSnap(snap string) (err error) {
	return d.UpdateSnap(d.appServerPath, snap)
}
func (d *ClusterClient) CloseAppServer() (err error) {
	return d.handler.Delete(d.appServerPath)
}
