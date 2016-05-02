package cluster

import "time"

func (d *spServer) ResetSPSnap() {
	services := d.services.GetService()
	for k := range services {
		nmap := d.getNewDataMap(k)
		path := nmap.Translate(serviceProviderPath)
		d.zkClient.ZkCli.UpdateValue(path, d.snap.GetSnap(k))
	}
}

func (d *spServer) StartRefreshSnap() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			d.ResetSPSnap()
		}
	}

}
