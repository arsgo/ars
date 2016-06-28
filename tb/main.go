package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "[Debug]", log.Llongfile)
}

//批量下单，指定进程数，创建指定进程的程序，并批量进行下单请求
func main() {

	var (
		totalRequest int
		concurrent   int
		timeout      int
		sleep        int
		configPath   string
	)

	flag.IntVar(&totalRequest, "n", 0, "总请求个数")
	flag.IntVar(&concurrent, "c", 1, "并发处理数")
	flag.IntVar(&timeout, "t", 0, "超时时长，默认不限制")
	flag.StringVar(&configPath, "f", "", "参数配置文件")
	flag.IntVar(&sleep, "s", 0, "每笔请求休息毫秒数")

	flag.Parse()
	config, err := NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	s, p := NewProcesss(totalRequest, concurrent, timeout, sleep, config)
	if !s {
		return
	}
	p.Close()
	calculateKPI(p.Start())

}
