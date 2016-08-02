package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/arsgo/lib4go/utility"
)

type conf struct {
	Cluster string `json:"cluster"`
	Domain  string `json:"domain"`
	Mask    string `json:"mask"`
}

//SysConfig 系统配置
type SysConfig struct {
	ZKServers []string
	Domain    string
	Mask      []string
	IP        string
}

var _config *SysConfig
var _err error
var _filePath string

func init() {
	_filePath = utility.GetExcPath("./conf/ars.conf.json", "bin")
	_config, _err = readConfig()
}

//SetConfig 设置服务器配置
func SetConfig(clusterServers string, domain string, mask string) {
	_config.ZKServers = strings.Split(clusterServers, ";")
	_config.Domain = domain
	_config.Mask = strings.Split(mask, ";")
}

//Get 获取配置
func Get() (*SysConfig, error) {
	return _config, _err
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
func createConfig(config *conf) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("创建配置文件", _filePath, "文件错误:", r, string(debug.Stack()))
		}
	}()
	data, _ := json.Marshal(config)
	ioutil.WriteFile(_filePath, data, os.ModeAppend)
}
func readConfig() (config *SysConfig, err error) {
	if !exist(_filePath) {
		fmt.Println("找不到配置文件:", _filePath)
		createConfig(&conf{})
	}

	configs := &conf{}
	bytes, err := ioutil.ReadFile(_filePath)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &configs); err != nil {
		fmt.Println("can't Unmarshal  ", _filePath, err.Error())
		return
	}

	if strings.EqualFold(configs.Cluster, "") || strings.EqualFold(configs.Domain, "") ||
		strings.EqualFold(configs.Mask, "") {
		err = fmt.Errorf("cluster:%s,domain:%s,mask:%s不能为空", configs.Cluster, configs.Domain,
			configs.Mask)
		return
	}
	config = &SysConfig{}
	config.ZKServers = strings.Split(configs.Cluster, ";")
	config.Domain = configs.Domain
	config.Mask = strings.Split(configs.Mask, ";")
	config.IP = utility.GetLocalIPAddress(config.Mask...)
	return config, err

}
func getDefConfig() *SysConfig {
	c := &SysConfig{}
	c.ZKServers = []string{"192.168.101.161:2181"} // []string{"171.221.206.81:2181"} //
	c.Domain = "/grs/weixin"
	c.Mask = []string{"192.168", "172.16"}
	c.IP = utility.GetLocalIPAddress(c.Mask...)
	return c
}
