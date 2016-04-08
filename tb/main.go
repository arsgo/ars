package main

import (
	"flag"
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
		address      string
		sleep        int
	)

	flag.IntVar(&totalRequest, "n", 0, "总请求个数")
	flag.IntVar(&concurrent, "c", 1, "并发处理数")
	flag.IntVar(&timeout, "t", 0, "超时时长，默认不限制")
	flag.StringVar(&address, "u", "", "请求的地址")
	flag.IntVar(&sleep, "s", 0, "每笔请求休息毫秒数")

	flag.Parse()

	s, p := NewProcesss(totalRequest, concurrent, address, timeout, sleep)
	if !s {
		return
	}

	response, totalMillisecond := p.Start()

	calculateKPI(response, totalMillisecond)

}
