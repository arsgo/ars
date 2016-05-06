
require 'xstring'
require 'xtable'


module ("xobject", package.seeall)

---检查对象中的字段是否为空
---1. luatable
---2. 字段名，多个用逗号分隔

empty = function(input, fields)
	local lst=xstring.split(fields,",")
	local max=xtable.size(lst)
	for i=1,max,1 do
		if(xstring.empty(input[lst[i]])) then
			return true
		end
	end
	return false
end

---根据输入的属性映射串拷贝对象，返回新的对象
---1. 输入对象luatable
---2. ”属性映射串“多个用逗号分隔,新旧属性名用">"分隔,属性名相同可省略">"和后面的部分，如: up_channel_no>channel_no,product_no
clone = function(input,content)
	local params={}
	local output={}
	for match in string.gmatch(tostring(content),"[^,]+") do
		local fname=match:match("([^>]+)")
		local sname=match:match(">([^>]+)")
		table.insert(params,{f=fname,s=sname or fname})
	end

	for i,v in pairs(params) do
		output[v.s]=input[v.f]
	end
	return output
end
