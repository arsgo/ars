
function main(input)

    local rpc=NewRPC()
	local session=rpc:AsyncRequest("get_pay_order","{}")
	r,e= rpc:GetAsyncResult(session)
	return r
end
