package cluster

import "time"

func (d *rcServer) StartRefreshSnap() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			d.ResetRCSnap()
		}
	}

}
