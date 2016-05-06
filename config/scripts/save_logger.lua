
function main(input)

    local rpc=NewRPC()
	local session=rpc:AsyncRequest("save_logger","{}")
	return rpc:GetAsyncResult(session)
end
