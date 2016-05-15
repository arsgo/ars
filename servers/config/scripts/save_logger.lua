
function main(input)

    local rpc=NewRPC()
	local session=rpc:AsyncRequest("save_logger","{}")
	r,e= rpc:GetAsyncResult(session)
	return r,e
end
