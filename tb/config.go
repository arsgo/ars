package main

import (
	"encoding/json"
	"io/ioutil"
)

type dataBlock struct {
	Params      json.RawMessage `json:"params"`
	CommandName string          `json:"command"`
}
type config struct {
	Items      []*dataBlock `json:"requests"`
	configPath string
	Address    string `json:"address"`
}

func NewConfig(configPath string) (*config, error) {
	return getRequestFromFile(configPath)
}

func getRequestFromFile(path string) (c *config, err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	c = &config{}
	if err = json.Unmarshal(bytes, &c); err != nil {
		return
	}
	return
}
