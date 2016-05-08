
require 'xstring'
require 'xnumber'
require 'xjson'
module ("xtable", package.seeall)
__xml_tag_config={
			XML_HEADER = "<?xml version=\"1.0\" encoding=\"%s\"?>",
	        XML_TARGET_BEGIN = "<%s>",
	        XML_TARGET_END = "</%s>",
	        XML_TARGET_LINE_BEGIN = "<%s",
	        XML_TARGET_LINE_END = " />",
	        XML_TARGET_ATTR_STR = " %s=\"%s\""
}

--- Return whether table is empty.
-- @param t table
-- @return <code>true</code> if empty or <code>false</code> otherwise
function empty (t)
  return not next (t)
end
--- Find the number of elements in a table.
-- @param t table
-- @return number of elements in t
size=function (t)
	if(not(t)) then
	 return 0
	end
  local n = 0
  for _ in pairs (t) do
    n = n + 1
  end
  return n
end

--- Make the list of indices of a table.
-- @param t table
-- @return list of indices
indices=function (t)
  local u = {}
  for i, v in pairs (t) do
    insert (u, i)
  end
  return u
end

--- Make the list of values of a table.
-- @param t table
-- @return list of values
values=function (t)
  local u = {}
  for i, v in pairs (t) do
    insert (u, v)
  end
  return u
end

--- Invert a table.
-- @param t table <code>{i=v, ...}</code>
-- @return inverted table <code>{v=i, ...}</code>
invert=function (t,k)
  local u = {}
  for i, v in pairs (t) do
	if(type(v)=="table") then
		u[v[k]]=v
	else
		u[v] = i
	end

  end
  return u
end


--- Make a shallow copy of a table, including any metatable (for a
-- deep copy, use tree.clone).
-- @param t table
-- @param nometa if non-nil don't copy metatable
-- @return copy of table
clone=function (t, nometa)
  local u = {}
  if not nometa then
    setmetatable (u, getmetatable (t))
  end
  for i, v in pairs (t) do
	if(type(v)=="table") then
		u[i] = clone(v)
	else
		u[i] = v
	end

  end
  return u
end

--- Merge two tables.
-- If there are duplicate fields, u's will be used. The metatable of
-- the returned table is that of t.
-- @param t first table
-- @param u second table
-- @return merged table
merge=function (t, u)
  local r = clone (t or {})
  for i, v in pairs (u) do
    r[i] = v
  end
  return r
end

concat=function(t,split,start)
	if(t==nil) then
		return ""
	end
	local index=start or 1
	local s=split or ""
	if(index==1) then
		return table.concat(t,s)
	else
		local r={}
		for i=index,#t,1 do
			table.insert(r,t[i])
		end
		return table.concat(r,s)
	end
end

append=function(t,u)
	for i, v in pairs (u) do
		table.insert(t,v)
	end
end



------------------------------------------------
---将TABLE拼接为以指定符号分隔的字符串
------------------------------------------------
--tb:数组或普通TABLE(需指定拼音的字段名)
--[splitChar]:分隔字符，缺省为","
--[fieldName]:字段名，为多字段的TABLE时需指定
------------------------------------------------
--示例:
--tb={}
--tb[1]={id="w"}
--tb[2]={id="f"}
--执行函数:sys.join(tb,",","id")
--输出结果:w,f
-------------------------------------------------
--举例2:
--tb={}
--tb[1]="w"
--tb[2]="w"
--sys.join(tb)
--输出结果:w,f
------------------------------------------------
join=function(tb,splitChar,fieldName)
	local spltChar=splitChar or ""
	local fname=fieldName or "*"
	if(not(isarray(tb))) then
		return nil
	end
	if(fname=="*") then
		return table.concat(tb,spltChar)
	end

	local rtb=keep(tb,fname)
	local ntb={}
	for i,v in pairs(rtb) do
		table.insert(ntb,v[fieldName])
	end
	return table.concat(ntb,spltChar)
end

-------------------------------------------------------------
---
-------------------------------------------------------------
keep=function(tb,f)
	if(type(f)=="number") then
		--保留指定条数的数据
		local v=tonumber(f)
		if(v>=xtable.size(tb)) then
			return tb
		else
			local rt={}
			for i=1,v,1 do
				rt[i]=tb[i]
			end
			return rt
		end
	else
	--保留字段，指定保留的字段，其它字段从table中移除
		local fieldName=f
		local fnames=xstring.split(fieldName,",")
		local names=invert(fnames)
		local rtb={}
		for i,v in pairs(tb) do
			rtb[i]={}
			for k,z in pairs(v) do
				if(names[k]~=nil) then
					rtb[i][k]=z
				end
			end
		end
		return rtb
	end
end

-------------------------------------------------------------
---移除字段，指定移除的字段，其它字段保留在table中
-------------------------------------------------------------
remove=function(tb,fieldNames)
	if(not(isarray(tb))) then
		return tb
	end
	local rtb={}
	local nfields=invert(xstring.split(fieldNames))
	for i,v in pairs(tb) do
		rtb[i]={}
		for k,d in pairs(v) do
			if(nfields[k]==nil) then
			  rtb[i][k]=d
			end
		end
	end
	return rtb
end
-------------------------------------------------------------
---是否是array数组，索引须从1开始
-------------------------------------------------------------
isarray=function(tb)
	if(type(tb)~="table") then
		return false
	end
	local n= xtable.size(tb)
	for i=1,n,1 do
	 if(tb[i]==nil) then
		return false
	 end
	end
	return true
end

__group=function(tb,nameFieldstb)
	local nameFields=nameFieldstb
	if(type(nameFieldstb)=="string") then
		local fnames=xstring.split(aname)
		nameFields=invert(fnames)
	end
	local maintb={}
	local othertb={}
	local keys={}
	for i,v in pairs(tb) do
		if(nameFields[i] or nameFields["*"] or xstring.start_with(i,tostring(nameFieldstb[1]))) then
			maintb[i]=v
			table.insert(keys,tostring(i).."-"..tostring(v))
		else
			othertb[i]=v
		end
	end
	table.sort(keys)
	return table.concat(keys),maintb,othertb
end

-------------------------------------------------------------
---分组数据
-------------------------------------------------------------
--示例:
--local data={}
--data[1]={type_id=1,type_name="接待率",value_id=1,value_name="100%"}
--data[2]={type_id=1,type_name="接待率",value_id=2,value_name="90%"}
--data[3]={type_id=2,type_name="转换率",value_id=3,value_name="99%"}
--print(xtype.to_json(group(data,"type_id,type_name","value_id,value_name")))
--输出结果为:
--[{"type_name":"转换率","type_id":2,"items":[{"value_id":3,"value_name":"99%"}]},
--{"type_name":"接待率","type_id":1,"items":[{"value_id":2,"value_name":"90%"},
--{"value_id":1,"value_name":"100%"}]}]
-------------------------------------------------------------
group=function(tb,aname,bname,cname,dname,ename,fname)
	---检查输入参数
	if(not(isarray(tb)) or aname==nil or tb ==nil) then
		return nil
	end
	---转换输入字段
	local rtb={}
	local sort_tb={}
	local fnames=xstring.split(aname)
	local names=invert(fnames)
	----分组规整数据
	for i,v in pairs(tb) do
		local key,m,o=__group(v,names)
		if(xtable.size(m)>0) then
			if(not(rtb[key])) then
				rtb[key]={}
				table.insert(sort_tb,rtb[key])
			end
			rtb[key].data=m
			rtb[key].items=rtb[key].items or {}
			table.insert(rtb[key].items,o)
		end
	end
	----规整子项
	local mtb={}
	for i,v in pairs(sort_tb) do
		local ntb={}
		ntb=v.data
		ntb.items=group(v.items,bname,cname,dname,ename,fname)
		table.insert(mtb,ntb)

	end
	return xtable.size(mtb)>0 and mtb or nil
end
print=function(tb)
	for k,v in pairs(tb) do
		_G.print(string.format("k:%s,v:%s",tostring(k),tostring(v)))
	end
end


sum=function(tb,f)
	local total=0
	for i,v in pairs(tb) do
		if(v[f]) then
			assert(tonumber(v[f]),string.format("索引:%s的值%s必须为数字",tostring(i),tostring(v[f])))
			total=total+tonumber(v[f])
		end
	end
	return total
end

max=function(tb,f)
	local max_value=0
	for i,v in pairs(tb) do
		if(v[f]) then
			local c=tonumber(v[f])
			assert(c,string.format("索引:%s的值%s必须为数字",tostring(i),tostring(c)))
			max_value=max_value>c and max_value or c
		end
	end
	return max_value
end

mul=function(tb,f,s)
	assert(tonumber(s),"arg3 输入参数必须为数字")
	for i,v in pairs(tb) do
		local c=v[f]
		if(c) then
			assert(c,string.format("索引:%s的值%s必须为数字",tostring(i),tostring(c)))
			v[f]=c*s
		end
	end
end



hasChildTable=function(tb)
	for i,v in pairs(tb) do
		if(type(v)=="table") then
			return true
		end
	end
	return false
end
rechange=function(source,names,start)
	local rtb={}
	local dmax=xtable.size(source)
	local nmax=xtable.size(names)
	local max=dmax>nmax and nmax or dmax
	for i=start,max,1 do
		rtb[names[i]]=source[i]
	end
	return rtb
end

tojson=function(t,c)
	assert(type(t)=="table","输入参数必须为table")
	if(xnumber.parse(c,0)==1) then
		return  xjson.encode(t[1])
	else
		return xjson.encode(t)
	end
end

toxml=function(tb,root,addheader,encoding,iselement)
	local root=root or "root"
	local addheader=addheader or false
	local encoding=encoding or "gb2312"
	local xml=addheader and string.format(__xml_tag_config.XML_HEADER,encoding) or ""
	local hasChild=hasChildTable(tb)
	local ise=hasChild or iselement
	xml=xml .. string.format(ise and __xml_tag_config.XML_TARGET_BEGIN or __xml_tag_config.XML_TARGET_LINE_BEGIN,
		root)

	if(isarray(tb)) then
		for k,v in pairs(tb) do
			if(type(v)=="table") then
				xml=xml..toxml(v,"item",false,encoding,iselement)
			end
		end
	else
		for k,v in pairs(tb) do
			if(type(v)=="table") then
				xml=xml..toxml(v,tostring(k),false,encoding,iselement)
			else
				if(ise) then
					xml=xml..string.format(__xml_tag_config.XML_TARGET_BEGIN,tostring(k))
					xml=xml..tostring(v)
					xml=xml..string.format(__xml_tag_config.XML_TARGET_END,tostring(k))
				else
					xml=xml..string.format(__xml_tag_config.XML_TARGET_ATTR_STR,tostring(k),tostring(v))
				end
			end
		end

	end
   xml=xml..string.format(ise and __xml_tag_config.XML_TARGET_END or __xml_tag_config.XML_TARGET_LINE_END,root)
   return xml
end



parse=function(t,i)
	assert(t,"输入参数不能为空")
	if(type(t)=="string") then
		assert(not(xstring.empty(t)),"输入参数不能为空")
		return xjson.decode(t)
	end

	assert(type(t)=="userdata" and type(t.get_row_count)=="userdata"
	and type(t.get_col_count)=="userdata" and type(t.get_col_name)=="userdata"
	and type(t.iget)=="userdata","输入参数必须为userdata,并实现接口get_row_count(),get_col_count(),get_col_name(i),iget(r,c)")
	local row_count=t:get_row_count()
	row_count=xnumber.min(i or row_count,row_count)
	local col_count=t:get_col_count()
	local r={}
	for i=1,row_count,1 do
		r[i]={}
		for j=1,col_count,1 do
			local name=t:get_col_name(j-1)
			r[i][name]=t:iget(i-1,j-1)
		end
	end
	return i==1 and r[1] or r
end

intercept = function(t, start, stop)
	if(#t < start) then
		return {}
	end
	local u = {}
	stop = stop > #t and #t or stop
	for i=start, stop, 1 do
		table.insert(u,t[i])
	end
	return u
end

