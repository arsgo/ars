
require 'xstring'

module ("xcache", package.seeall)
-------------------String Start------------------------
function xcache:get(key)
	local data = CacheFactory.Get():Get(key)
	if(xstring.empty(data)) then
		return ''
	end
	return data
end

function xcache:set(key, value, seconds)
	if(type(value) == 'table') then
		value = xtable.tojson(value)
	end
	CacheFactory.Get():Set(key, value, seconds)
end
