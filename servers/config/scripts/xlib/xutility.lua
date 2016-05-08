
local string=string
local tostring=tostring
local utility=Utility
module ("xutility")

decode=function(s)
	return string.gsub(s, '%%(%x%x)', function(h) return string.char(tonumber(h, 16)) end)
end
encode=function(s)
	local  v= string.gsub(s, "([^%w%.%- ])", function(c) return string.format("%%%02X", string.byte(c)) end)
    return string.gsub(v, " ", "+")
end

md5=function(s,c)
	return utility.md5(s, c or "utf-8")
end








