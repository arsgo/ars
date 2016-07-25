package base

import "github.com/colinyl/lib4go/pool"

type ProxySnap struct {
	pool.ObjectSnap
	ElapsedTime ServerSnap `json:"elapsed"`
}

//getProxySnap 获取当前脚本
func GetProxySnap(objectPoolSnaps pool.ObjectPoolSnap, snaps map[string]interface{}) (r []interface{}) {
	poolSnaps := objectPoolSnaps.Snaps
	if len(poolSnaps) == 0 {
		r = make([]interface{}, 0)
		return
	}
	for _, v := range poolSnaps {
		if elp, ok := snaps[v.Name]; ok {
			sr := elp.(*ProxySnap)
			sr.Cache = v.Cache
			sr.MinSize = v.MinSize
			sr.MaxSize = v.MaxSize
			sr.Name = v.Name
			sr.Status = v.Status
			//	sr.Created = v.Created
			r = append(r, sr)
		} else {
			r = append(r, v)
		}
	}
	return r
}
