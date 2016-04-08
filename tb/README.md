# http bench
        go get -u -v github.com/colinyl/hb
        
###参数说明
 -n int       总请求个数
 
 -c int       并发请求数
 
 -t int       超时时长，默认不限制
 
 -s int      每笔请求休息毫秒数
 
 -u string    请求的URL,未指定时需通过-f参数指定参数文件
 
 -f string    参数配置文件,未指定时需使用-u指定URL
       

 
       
###示例

         > hb -n 100 -c 100 -f config_jyk.json
         启动 100 个工作进程,处理 100个请求

         -------------------------------------------------------------------------
         总数    成功    平均耗时        每秒请求数
         100     100     545.37          109.29
         -------------------------------------------------------------------------
         
###配置文件
        [
            {
                "params": {
                    "orderNo": "OR@guid",
                    "coopId": "0",
                    "productStandard": "100",
                    "orderNum": "1",
                    "totalStandard": "100",
                    "rechargeAccount": "",
                    "notifyUrl": "http://192.168.101.139:9999/t/order/notify",
                    "timestamp": "@timestamp",
                    "$": "38b8b230cbab4459987073990089c8ae{@raw}38b8b230cbab4459987073990089c8ae"
                },
                "url": "http://192.168.101.139:9999/order/request"
         }
     ]






