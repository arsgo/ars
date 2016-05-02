package cluster

import (
	"fmt"

	"github.com/colinyl/ars/rpcservice"
	"github.com/colinyl/lib4go/logger"
)

func (a *appServer) StopJobServer() {
	if a.jobServer != nil {
		a.jobServer.Stop()
	}
}
func (a *appServer) StartJobConsumer(jobNames []string) {
	a.lk.Lock()
	defer a.lk.Unlock()
	if !a.hasStartJobServer && len(jobNames) > 0 {
		a.Log.Info("start rpc service for job consumer")
		a.hasStartJobServer = true
		a.jobServerAdress = rpcservice.GetLocalRandomAddress()
		a.snap.Address = fmt.Sprintf("%s%s", a.zkClient.LocalIP, a.jobServerAdress)
		a.jobServer = rpcservice.NewRPCServer(a.jobServerAdress, &appServerJobHandler{server: a, Log: a.Log})
		err := a.jobServer.Serve()
		if err != nil {
			a.Log.Error(err)
		}
	}
	jobMap := make(map[string]string)
	for _, v := range jobNames {
		jobMap[v] = v
	}
	//clear
	for i := range a.jobNames {
		if p, ok := jobMap[i]; !ok {
			a.zkClient.ZkCli.Delete(p)
		}
	}
	//add jobConsumerPath   = "@domain/job/@jobName/consumers/job_"
	dmap := a.dataMap.Copy()
	for i := range jobMap {
		if _, ok := a.jobNames[i]; !ok {
			dmap.Set("jobName", i)
			path, err := a.zkClient.ZkCli.CreateSeqNode(dmap.Translate(jobConsumerPath),
				a.snap.GetSnap())
			if err != nil {
				a.Log.Error(err)
				continue
			}
			a.jobNames[i] = path
			a.Log.Infof("::start job service:%s", i)
		}
	}

}

type appServerJobHandler struct {
	server *appServer
	Log    *logger.Logger
}

func (r *appServerJobHandler) Request(name string, input string) (result string, err error) {
	r.Log.Info("recv job server request")
	return getSuccessResult(), nil
}
func (r *appServerJobHandler) Send(name string, input string, data []byte) (string, error) {
	return getErrorResult("200", "not support send methods"), nil
}
func (r *appServerJobHandler) Get(name string, input string) ([]byte, error) {
	return []byte{}, nil
}
