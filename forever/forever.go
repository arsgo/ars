package forever

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/colinyl/daemon"
	"github.com/colinyl/lib4go/logger"
)

type forever struct {
	dm   daemon.Daemon
	log  *logger.Logger
	svs  service
	name string
	desc string
}
type service interface {
	Start() error
	Stop() error
}

func NewForever(svs service, log *logger.Logger, name string, desc string) *forever {
	dm, err := daemon.New(name, desc)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &forever{dm: dm, name: name, desc: desc, svs: svs, log: log}
}
func (f *forever) Start() {
	defer func(){
		if r:=recover();r!=nil{
			f.log.Error(r)
			fmt.Println("error happend is write to file")
		}
	}()
	result, err := f.run()
	if err != nil {
		f.log.Error(err)
		return
	}
	f.log.Info(result)
}

func (f *forever) run() (string, error) {

	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", f.name)
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return f.dm.Install()
		case "remove":
			return f.dm.Remove()
		case "start":
			return f.dm.Start()
		case "stop":
			return f.dm.Stop()
		case "status":
			return f.dm.Status()
		default:
			return usage, nil
		}
	}
	if err := f.svs.Start(); err != nil {
		return "", err
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	for {
		select {
		case <-interrupt:
			f.svs.Stop()
			return fmt.Sprintf("%s was killed", f.name), nil
		}
	}
	// never happen, but need to complete code
	return usage, nil
}
