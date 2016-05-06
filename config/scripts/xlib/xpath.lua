
require 'xtable'
module("xpath", package.seeall)

xpath.get_name=function(content)
	local pattern = "([^%[%]]+)[%[]*"
	local match = string.match(tostring(content),pattern)
	return match~=nil,match
end
xpath.get_params_key=function(content)
	local pattern = "([^=,@%]]+)"
	local match = string.match(content,pattern)
	return match~=nil,match
end

xpath.get_params_value=function(content)
	local pattern = "=([^,%]]+)"
	local match = string.match(content,pattern)
	return match~=nil,match
end
xpath.get_params=function(content)
	local mode={}
	local s,n=xpath.get_name(content)
	mode.tag=n
	local pattern = "%[([^%]]+)%]"
	local match = string.match(tostring(content),pattern)
	if(match==nil) then
		return mode
	end
	mode.attrs={}
	local items=xstring.split(match,",")
	for i,v in pairs(items) do
		i=tonumber(v)
		if(i) then
			mode.index=i
		else
			local s,key=xpath.get_params_key(v)
			local s,value=xpath.get_params_value(v)
			mode.attrs[key]=value
		end
	end
	return mode
end
xpath.parse=function(path)
	local r={}
	local pattern = "([^/]+)"
	for i,v in string.gmatch(tostring("/"..path.."/"),pattern) do
		table.insert(r,xpath.get_params(i))
	end
	return r
end





