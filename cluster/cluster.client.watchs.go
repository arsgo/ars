package cluster

import "sync"

type ValueChangeWatcher struct {
	callback      func(path string)
	paths         map[string]string
	clusterClient IClusterClient
	lock          sync.Mutex
}

func NewValueChangeWatcher(client IClusterClient, callback func(path string)) *ValueChangeWatcher {
	wp := &ValueChangeWatcher{clusterClient: client, callback: callback}
	wp.paths = make(map[string]string)
	return wp
}

func (w *ValueChangeWatcher) Push(path string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	if _, ok := w.paths[path]; ok {
		return
	}
	w.paths[path] = path
	w.clusterClient.WatchClusterValueChange(path, func() {
		w.callback(path)
	})
}
