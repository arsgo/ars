module ("xnumber", package.seeall)

-----------------------------------
---��������ַ����Ƿ�������
---@s,�ַ���
-----------------------------------
check=function (...)
	local arg=_VERSION=="Lua 5.1" and arg or  {...}
	if(#arg==0) then
		return false
	end
	for i=1,#arg do
		if(tonumber(arg[i])==nil) then
			return false
		end
	end
  return true
end


min=function(...)

	local arg=_VERSION=="Lua 5.1" and arg or  {...}
	assert(#arg~=0,"�����������Ϊ��")
	local min_value
	for i=1,#arg do
		local v=tonumber(arg[i])
		assert(v,"�����������Ϊ����")
		min_value=(min_value==nil or v<min_value) and v or min_value
	end
	return min_value
end


max=function(...)

	local arg=_VERSION=="Lua 5.1" and arg or  {...}
	assert(#arg~=0,"�����������Ϊ��")
	local max_value
	for i=1,#arg do
		local v=tonumber(arg[i])
		assert(v,"�����������Ϊ����")
		max_value=(max_value==nil or v>max_value) and v or max_value
	end
	return max_value
end

function parse(v,d)
	local r=tonumber(v)
	if(not(r)) then
		return d
	end
	return r
end

