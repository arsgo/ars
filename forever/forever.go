package forever

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/colinyl/daemon"
)

type forever struct {
	dm   daemon.Daemon
	name string
	desc string
}

func NewForever(name string, desc string) *forever {
	dm, err := daemon.New(name, desc,"log","github.com/colinyl/lib4go/logger")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &forever{dm: dm, name: name, desc: desc}
}

func (f *forever) Manage(start func()(interface{}), close func(o interface{})) (string, error) {

	usage := fmt.Sprintf("Usage: %s install | remove | start | stop | status", f.name)
	// if received any kind of command, do it
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

	obj:=start()

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case <-interrupt:
			close(obj)
			return fmt.Sprintf("%s was killed",f.name), nil
		}
	}
	// never happen, but need to complete code
	return usage, nil
}
