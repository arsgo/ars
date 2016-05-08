
function main(input)

    local rpc=NewRPC()
	local session=rpc:AsyncRequest("get_pay_order","{}")
	return rpc:GetAsyncResult(session)
end
