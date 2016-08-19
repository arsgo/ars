package main

import (
	"testing"

	"github.com/arsgo/ars/servers/config"
)

func TestJobCall(t *testing.T) {
	conf := config.GetDefConfig()
	t.Logf("conf:%+v", conf)
	rcserver, err := NewRCServer(conf)
	if err != nil {
		t.Error(err)
	}
	err = rcserver.Start()
	if err != nil {
		t.Error(err)
	}

}
