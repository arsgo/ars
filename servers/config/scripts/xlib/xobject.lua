
require 'xstring'
require 'xtable'


module ("xobject", package.seeall)

---�������е��ֶ��Ƿ�Ϊ��
---1. luatable
---2. �ֶ���������ö��ŷָ�

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

---�������������ӳ�䴮�������󣬷����µĶ���
---1. �������luatable
---2. ������ӳ�䴮������ö��ŷָ�,�¾���������">"�ָ�,��������ͬ��ʡ��">"�ͺ���Ĳ��֣���: up_channel_no>channel_no,product_no
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
