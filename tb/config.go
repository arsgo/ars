package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)
type dataBlock struct {
	Params    map[string]string
	URL string
}
type config struct {
	Items      []*dataBlock
	configPath string
	URL        string
}

func NewConfig(configPath string, URL string) *config {
	c := &config{configPath: configPath, URL: URL}
	c.getRequestData()
	return c
}

func (c *config) getRequestFromFile() {
	bytes, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.Items = make([]*dataBlock, 0)
	if err := json.Unmarshal(bytes, &c.Items); err != nil {
		fmt.Printf("72:can't Unmarshal %s:%s\r\n", c.configPath, err.Error())
		os.Exit(1)
	}
    for i:=0;i<len(c.Items);i++{
        if  strings.EqualFold(c.Items[i].URL,""){
            fmt.Printf("配置:%d URL不能为空\r\n",i)
            os.Exit(1)
        }
    }
}

func (c *config) getRequestData() {
	if !strings.EqualFold(c.configPath, "") {
		c.getRequestFromFile()
        return
	}
	c.Items = make([]*dataBlock, 0)
	c.Items = append(c.Items, &dataBlock{})
	c.Items[0].Params = make(map[string]string, 0)
	c.Items[0].URL = c.URL
}
