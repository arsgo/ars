package cluster

func (d *ClusterClient) CreateNode(path string, value string) (p string, err error) {
	p, err = d.handler.CreateTmpNode(path, value)
	return
}
func (d *ClusterClient) UpdateNode(path string, value string) (err error) {
	err = d.handler.UpdateValue(path, value)
	return
}
func (d *ClusterClient) CloseNode(path string) (err error) {
	return d.handler.Delete(path)
}
func (d *ClusterClient) SetNode(path string, value string) (err error) {
	if _, ok := d.handler.Exists(path); !ok {
		_, err = d.CreateNode(path, value)
	} else {
		err = d.UpdateNode(path, value)
	}
	return
}
