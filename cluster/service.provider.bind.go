package cluster

import (
	"errors"
	"strings"
	"time"
)

func (d *spServer) WatchServiceConfigChange() {
	d.zkClient.waitZKPathExists(d.serviceConfig, time.Hour*8640, func(exists bool) {
		if !exists {
			d.Log.Info("sp server config not exists")

		} else {
			d.rebind()
		}
	})
	d.Log.Info("::watch for provider config change")
	d.zkClient.watchZKValueChange(d.serviceConfig, func() {
		d.rebind()
	})
}

func (d *spServer) rebind() {
	d.lk.Lock()
	defer d.lk.Unlock()
	defer func() {
		d.Log.Infof("bind services:%s,%s", d.mode, d.services.ToString())
	}()
	aloneService, sharedService := d.groupService()
	if strings.EqualFold(d.mode, eModeAlone) {
		goon, _ := d.checkAloneService(aloneService)
		if !goon {
			return
		}
	}
	config, _ := d.bindServices(aloneService)
	defer d.setMQConsumer(d.services.services)
	if len(config.services) > 0 {
		d.deleteSharedSevices(config.services)
		d.mode = eModeAlone
		d.services.Reset(config.services)
		return
	}

	config, _ = d.bindServices(sharedService)
	d.mode = eModeShared
	d.deleteSharedSevices(config.services)
	d.services.Reset(config.services)
}

func (d *spServer) deleteSharedSevices(svs map[string]spService) {
	for i := range d.services.services {
		if _, ok := svs[i]; ok {
			continue
		}
		nmap := d.getNewDataMap(i)
		path := nmap.Translate(serviceProviderPath)
		d.zkClient.ZkCli.Delete(path)
	}
}

func (d *spServer) checkAloneService(configs map[string]spService) (ct bool, err error) {
	ct = true
	err = nil
	if len(d.services.services) < 1 {
		return
	} else if len(d.services.services) > 1 {
		for i := range d.services.services {
			nmap := d.getNewDataMap(i)
			d.deleteSPPath(nmap.Translate(serviceProviderPath))
		}
		return
	}
	for i, v := range d.services.services {
		if cc, ok := configs[i]; ok && (strings.EqualFold(v.getUNIQ(), cc.getUNIQ()) || d.zkClient.checkIP(configs[i].IP)) {
			nmap := d.getNewDataMap(i)
			path := nmap.Translate(serviceProviderPath)
			if d.zkClient.ZkCli.Exists(path) {
				ct = false
				return
			}
		}
	}
	return
}

func (d *spServer) bindServices(services map[string]spService) (psconfig *spConfig, err error) {
	psconfig = &spConfig{services: make(map[string]spService, 0)}
	for sv, config := range services {
		err = nil
		if v, ok := d.services.services[sv]; !ok {
			err = d.bindService(sv, config)
		} else if !strings.EqualFold(config.getUNIQ(), v.getUNIQ()) {
			err = d.bindService(sv, config)
		}
		if err == nil {
			psconfig.services[sv] = config
		}
		if strings.EqualFold(config.Mode, eModeAlone) && err == nil {
			return
		}
	}
	return
}

func (d *spServer) bindService(serviceName string, config spService) (err error) {
	nmap := d.getNewDataMap(serviceName)
	path := nmap.Translate(serviceProviderPath)
	if !d.zkClient.checkIP(config.IP) {
		d.deleteSPPath(path)
		return errors.New("ip not match")
	}
	d.createSPPath(path, d.snap.GetSnap(serviceName))
	return
}
