package cluster

import "time"

func (d *appServer) StartRefreshSnap() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			d.ResetAPPSnap()
		}
	}

}
