

local xstring=require('xstring')
local xtable = require("xtable")
local xutility = require("xutility")
local xdate=require('xdate')
local base = _G
local table=table
local print=print
local string=string
local tostring=tostring


module ("xsign")

check=function(input,token)
	local sign=input.sign
	local csign,raw=get_sign(input,token)
	print("csign"..csign)
	return string.upper(sign)==string.upper(csign),raw
end

get_sign=function(input,token)
	local raw=token..get_raw(input)..token
	print("Ç©ÃûÔ­´®£º"..raw)
	return xutility.md5(raw),raw
end

get_raw=function(input)
	local qinput=xtable.clone(input)
	print(qinput)
	qinput.sign=nil
	local sort_input={}
	for i,v in base.pairs(qinput) do
		table.insert(sort_input,string.format("%s%s",tostring(i),tostring(v)))
	end
	table.sort(sort_input)
	local raw=table.concat(sort_input)
	return raw
end


expire=function(timestamp,mi)
	return xdate:now()>xdate:new(timestamp,"yyyyMMddhhmmss"):addminutes(5)
end
