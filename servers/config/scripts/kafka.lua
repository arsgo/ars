sys = require("sys")
--rpc=require("rpc")

function main(input)
	kafka=mqp.new("kafka01")
    result=kafka:publish("get_pay_order")	
	return result
end
