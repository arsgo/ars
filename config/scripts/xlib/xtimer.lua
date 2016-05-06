

xtimer={timeout=3}

xtimer.set=function(timeout)
	xtimer.timeout=timeout
end
xtimer.sleep=function(n)
	os.execute("sleep "..n)
end
xtimer.start=function()
	if(xtimer.start_time==nil) then
		xtimer.start_time=os.time()
	end
end
xtimer.reset=function()
	xtimer.start_time=os.time()
end
xtimer.stop=function()
	xtimer.start_time=nil
end
xtimer.istimeout=function()
	if(xtimer.timeout==0) then
		return false
	end
	xtimer.start()
	return os.time()-xtimer.start_time>xtimer.timeout
end




