package cluster

import "time"

func (d *rcServer) StartSnapValue() {
	tp := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-tp.C:
			d.setLastParams()
			d.resetRCSnap()
		}
	}
}
