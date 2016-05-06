require("xlib/xjson")
function main(input)

    local config=xjson.decode(input)
    local elastic,err=NewElastic(config.params.elastic)
	  local id,err=elastic:Create("logger","delivery",[[{"name":"coliny_hh","content":"I don't know why because the build and tests a"}]])
    print("save:",id,err)
    print("query:",elastic:Search("logger","delivery",[[{"query":{"term":{ "name":"coliny_hh" }}}]]))
    local mq,err=NewMQProducer(config.params.mq)
    print("mq:",mq,err)
    
    err=mq:Send(config.params.queue,input)
     print("mq.send:",err)
    
	--return result,err
end
