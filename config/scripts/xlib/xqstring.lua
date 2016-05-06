require 'xstring'
require 'xtable'

module ("xqstring", package.seeall)
function xqstring:new(v)
	local o = {}
     setmetatable(o,self)
     self.__index = self
	 o.params={}
	 local iv=v or {}
	 if(type(iv)=="table") then
		if(iv.params) then
		  o.params=iv.params or {}
		else
			for k,v in pairs(iv) do
			 table.insert(o.params,{k=k,v=tostring(v)})
			end
		end

	 elseif(type(iv)=="userdata") then
		assert(type(iv.get_value)=="userdata" and type(iv.get_keys)=="userdata","输入参数未实现接get_keys()或get_value(key)")
		local keys=iv:get_keys()
		local tkey=xstring.split(keys)
		for i=1,#tkey,1 do
			local k=tkey[i]
			local v=iv:get_value(k)
			table.insert(o.params,{k=k,v=v})
		end
	 else
		assert(false,"输入参数必须是table或userdata")
	 end
     o.config={kvc="=",sc="&",req=false,ckey=true,nullkey=false,encoding="gbk"}
     return o
end


function xqstring:decode(encoding)
	local e =encoding or "gbk"
	local nparams={}
	for i,v in pairs(self.params) do
		local key=xstring.convert("gbk",e,v.k)
		local value=xstring.convert("gbk",e,v.v)
		print(string.format("key:o:[%s],c:[%s]",v.k,url.d(key,e)))
		print(string.format("value:o:[%s],c:[%s]",v.v,url.d(value,e)))

		table.insert(nparams,{k=url.d(key,e),v=url.d(value,e)})
	end
	self.params=nparams
end




function parse(str, eq, sep)
  if not sep then sep = '&' end
  if not eq then eq = '=' end
 local qstr=xqstring:new()
  for pair in string.gmatch(tostring(str), '[^' .. sep .. ']+') do
    if not string.find(pair, eq) then
    	qstr:add(pair,"")
    else
      local key, value = string.match(pair, '([^' .. eq .. ']*)' .. eq .. '(.*)')
      if key then
      	qstr:add(key,value)
      end
    end
  end
  return qstr
end



function xqstring:add(k,v)
	table.insert(self.params,{k=k,v=tostring(v)})
end

function xqstring:find(k)
	local s=xtable.size(self.params)
	for i=1,s,1 do
		if(self.params[i].k==k) then
			return i
		end
	end
	return -1
end
function xqstring:get(k)
	local index=self:find(k)
	if(index<1) then
		return nil
	end
	return self.params[index].v
end

function xqstring:get_all()
	return self.params
end

function xqstring:remove(k)
	local index=self:find(k)
	if(index>0) then
		table.remove(self.params,index)
	end

end
function xqstring:clear()
	self.params={}
end

function xqstring:sort()
	table.sort(self.params,function(x,y)
		return x.k<y.k
	end)
end

function xqstring:make(cf)
	local rconfig=cf or {}
	rconfig=xtable.merge(self.config,rconfig)
	local str=""
	for i,v in pairs(self.params) do
		if(not(rconfig.req) or v.v~=nil) then
			str=str..string.format("%s%s%s%s",
				(rconfig.ckey or (v.v==nil and nullkey)) and v.k or "",
				(rconfig.ckey or (v.v==nil and nullkey)) and rconfig.kvc or "",
				v.v,rconfig.sc)
		end
	end
	return xstring.trim(str,rconfig.sc)
end


