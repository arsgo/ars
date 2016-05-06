
local sys_totstring=tostring
module ("xstring", package.seeall)

_COPYRIGHT = "Copyright (C) 2014-2024"
_DESCRIPTION = "xstring"
_VERSION = "1.0"


-----------------------------------
---从字符串移除前导匹配项
---@s,字符串
---@r,正则表达式
-----------------------------------
ltrim=function (s, r)
  r = r or "%s+"
  return (string.gsub (s, "^" .. r, ""))
end

-----------------------------------
---从字符串移除尾部匹配项
---@s,字符串
---@r,正则表达式
-----------------------------------
rtrim=function (s, r)
  r = r or "%s+"
  return (string.gsub (s, r .. "$", ""))
end

-----------------------------------
---从字符串移除前导匹配项和尾部匹配项
---@s,字符串
---@r,正则表达式
-----------------------------------
trim =function  (s, r)
  return rtrim (ltrim (s, r), r)
end


function tfind (s, p, init, plain)
  local function pack (from, to, ...)
    return from, to, {...}
  end
  return pack (p.find (s, p, init, plain))
end
function finds (s, p, init, plain)
  init = init or 1
  local l = {}
  local from, to, r
  repeat
    from, to, r = tfind (s, p, init, plain)
    if from ~= nil then
      table.insert (l, {from, to, capt = r})
      init = to + 1
    end
  until not from
  return l
end

function concat (...)
  local r = {}
  for _, l in ipairs ({...}) do
    for _, v in ipairs (l) do
      table.insert (r, v)
    end
  end
  return r
end
function flatten (l)
  local m = {}
  for _, v in ipairs (l) do
    if type (v) == "table" then
      m = concat (m, flatten (v))
    else
      table.insert (m, v)
    end
  end
  return m
end
-----------------------------------
---将字符串以指字分隔符截取，并返回截断后的数组
---@s,字符串
---@r,分隔符
-----------------------------------
split=function (s, p)
	local sep=p or ","
	local pairs = concat ({0}, flatten (finds (s, sep)), {0})
	local l = {}
	for i = 1, #pairs, 2 do
		table.insert (l, string.sub (s, pairs[i] + 1, pairs[i + 1] - 1))
	end
	return l
end

-----------------------------------
---返回字符串是否为空
---@s,字符串
-----------------------------------
empty=function (...)
	local arg=_VERSION=="Lua 5.1" and arg or {...}
	if(#arg==0) then
		return true
	end
	for i=1,#arg do
		if(arg[i]==nil or tostring(arg[i])=="") then
			return true
		end
	end
  return false
end


-----------------------------------
---检查字符串是否以指定字符开头
---@param,t 待检查字符串
---@param,v 匹配字符串
-----------------------------------
start_with=function(t,v)
	return string.sub(t,1,string.len(v))==v
end

end_with=function(t,v)
	local tn=string.len(t)
	local vn=string.len(v)
	return (tn>vn) and string.sub(t,tn-vn+1)==v or false
end

encode=function(t,c)
 return t
end

__format=function(s,t,c)
	local pattern = "({[@#][%w_.]+})"
	local dst = s
	local input=t or {}
	local encoding=c or "gbk"
	for match in string.gmatch(dst,pattern) do
        local name = string.sub(match,3,string.len(match)-1)
        local needencode = false
		if(string.sub(match,2,2) == "#") then
		    needencode = true
		end

        if(input[name] == nil) then
			dst = string.gsub(dst,match,"")
		elseif(needencode) then
		    local v = xstring.encode(tostring(input[name]),encoding)
			v = string.gsub(v,"%%","%%%%")
			dst = string.gsub(dst,match,v)
		else
            local v = tostring(input[name])
			v = string.gsub(v,"%%","%%%%")
            dst = string.gsub(dst,match,v)
		end
    end
    return dst
end

format=function (f, arg1, arg2,...)
	local arg=_VERSION=="Lua 5.1" and arg or {...}
	if arg1 == nil then
		return f
	end
	if(type(arg1)=="table") then
		local xformat= __format(f,arg1,arg2)
		if(arg~=nil) then
			return string.format(xformat,unpack(arg))
		end
		return xformat
	else
		return string.format(f, arg1,arg2,...)
	end
end

parse=function(s)
	return s==nil and "" or tostring(s)
end

rpt=function(l,c)
	local t=""
	local count=l/string.len(c)
	while(count>0) do
		t=t..c
		count=count-1
	end
	return string.len(t)==l and t or string.sub(t..c,1,l)
end

lpading=function(s,l,c)
	return rpt(l-string.len(s),c)..s
end

rpading=function(s,l,c)
	return s..rpt(l-string.len(s),c)
end


join=function( ... )
	local r=""
	local s=","
	local arg=_VERSION=="Lua 5.1" and arg or {...}
	for i,x in ipairs(arg) do
		if(type(x)=="table") then
			for k,v in pairs(x) do
				r=r..s..tostring(v)
			end
		else
			r=r..s..tostring(x)
		end
   end
   return trim(r,s)
end

convert=function(fe,te,s)
	return utility.convert(fe,te,s)
end

tostring=function(s)
 return s==nil and "" or sys_totstring(s)
end

equals=function(a,b)
	return tostring(a) == tostring(b)
end
