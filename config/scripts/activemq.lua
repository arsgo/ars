sys = require("sys")
--rpc=require("rpc")

function main(input)
   print("start....................")
	m=MQ.new("stomp")
	print("send")
    result=m:send("go:stomp:test","123456")	
	print("consumer")
	
--	m:consume("go:stomp:test",callback)
	m:close()
	return result
end
function callback(msg)
	print(msg)
end