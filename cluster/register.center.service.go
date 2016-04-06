package cluster

import (
	"encoding/json"
	"time"
)

func (d *rcServer) WatchServiceChange(callback func(services map[string][]string, err error)){
    waitZKPathExists(d.servicePublishPath,time.Hour*8640,func(exists bool){
        if !exists {
			d.Log.Info("service publish config not exists")
		} else {          
			callback(d.getServiceProviderRoot())
		}
    })
    d.Log.Info("::watch for service config changes ")
    watchZKValueChange(d.servicePublishPath,func(){
        callback(d.getServiceProviderRoot())
    })
}


//WatchServiceProviderChange watch whether any service privider is changed
func (d *rcServer) watchServiceProviderChange() (err error) {
	if !d.IsMasterServer {
		return
	}
	d.Log.Info("::watch for service providers changes")
	waitZKPathExists(d.serviceRoot, time.Hour*8640, func(exists bool) {
		if !exists {
			d.Log.Info("service provider node not exists")
		} else {
			err = d.serviceChange()
		}
	})

	watchZKChildrenPathChange(d.serviceRoot, func() {
		d.serviceChange()
	})
	sproots, err := d.getServiceProviderRoot()
	for _, v := range sproots {
		for _, p := range v {
			watchZKChildrenPathChange(p, func() {
				err = d.serviceChange()
			})
		}

	}
    return
}

func (d *rcServer) getServiceProviderRoot() (map[string][]string, error) {
	var spList ServiceProviderList = make(map[string][]string)
	serviceList, err := zkClient.ZkCli.GetChildren(d.dataMap.Translate(serviceRoot))
	if err != nil {
		return spList, err
	}

	for _, v := range serviceList {
		nmap := d.dataMap.Copy()
		nmap.Set("serviceName", v)
		providerList, er := zkClient.ZkCli.GetChildren(nmap.Translate(serviceProviderRoot))
		if er != nil {
			return spList, er
		}
		for _, l := range providerList {
			spList.Add(v, l)
		}
	}
	return spList, nil
}

func (d *rcServer) getSPServices() (ServiceProviderList, error) {
	var spList ServiceProviderList = make(map[string][]string)
	serviceList, err := zkClient.ZkCli.GetChildren(d.dataMap.Translate(serviceRoot))
	if err != nil {
		return spList, err
	}

	for _, v := range serviceList {
		nmap := d.dataMap.Copy()
		nmap.Set("serviceName", v)
		providerList, er := zkClient.ZkCli.GetChildren(nmap.Translate(serviceProviderRoot))
		if er != nil {
			return spList, er
		}
		for _, l := range providerList {
			spList.Add(v, l)
		}
	}
	return spList, nil
}

func (d *rcServer) serviceChange() (err error) {
	d.setServiceParams()
	d.resetRCSnap()
	err = d.publishServices()
	return
}

func (d *rcServer) publishServices() (err error) {
	if !d.IsMasterServer {
		return
	}
	providers, err := d.getSPServices()
	if err != nil {
		return
	}
	buffer, err := json.Marshal(providers)
	if err != nil {
		return
	}
	serviceValue := string(buffer)
	path := d.dataMap.Translate(d.servicePublishPath)
	if zkClient.ZkCli.Exists(path) {
		err = zkClient.ZkCli.UpdateValue(path, serviceValue)
	} else {
		err = zkClient.ZkCli.CreatePath(path, serviceValue)
	}
	return
}
