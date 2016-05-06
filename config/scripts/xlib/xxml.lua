require 'lib4net'
require 'xtable'
require 'xpath'

module("xxml", package.seeall)

---------------------------------------------------------------------------------
---------------------------------------------------------------------------------
--
-- xml.lua - XML parser for use with the Corona SDK.
--
-- version: 1.2
--
-- CHANGELOG:
--
-- 1.2 - Created new structure for returned table
-- 1.1 - Fixed base directory issue with the loadFile() function.
--
-- NOTE: This is a modified version of Alexander Makeev's Lua-only XML parser
-- found here: http://lua-users.org/wiki/LuaXml
--
---------------------------------------------------------------------------------
---------------------------------------------------------------------------------




    function xxml:toxmlstring(value)
        value = string.gsub(value, "&", "&amp;"); -- '&' -> "&amp;"
        value = string.gsub(value, "<", "&lt;"); -- '<' -> "&lt;"
        value = string.gsub(value, ">", "&gt;"); -- '>' -> "&gt;"
        value = string.gsub(value, "\"", "&quot;"); -- '"' -> "&quot;"
        value = string.gsub(value, "([^%w%&%;%p%\t% ])",
            function(c)
                return string.format("&#x%X;", string.byte(c))
            end);
        return value;
    end

    function xxml:fromxmlstring(value)
        value = string.gsub(value, "&#x([%x]+)%;",
            function(h)
                return string.char(tonumber(h, 16))
            end);
        value = string.gsub(value, "&#([0-9]+)%;",
            function(h)
                return string.char(tonumber(h, 10))
            end);
        value = string.gsub(value, "&quot;", "\"");
        value = string.gsub(value, "&apos;", "'");
        value = string.gsub(value, "&gt;", ">");
        value = string.gsub(value, "&lt;", "<");
        value = string.gsub(value, "&amp;", "&");
        return value;
    end

    function xxml:ParseArgs(node, s)
        string.gsub(s, "(%w+)=([\"'])(.-)%2", function(w, _, a)
            node:addProperty(w, self:fromxmlstring(a))
        end)
    end

    function xxml:load(xmlText)
        local stack = {}
        local top = newNode()
        table.insert(stack, top)
        local ni, c, label, xarg, empty
        local i, j = 1, 1
		local cdata_len = 0
        while true do
            ni, j, c, label, xarg, empty = string.find(xmlText, "<(%/?)([%w_:]+)(.-)(%/?)>", i)
            if not ni then break end
            local text = string.sub(xmlText, i, ni - 1)
			text, cdata_len = string.gsub(text, "^<!%[CDATA%[", "")
			if(cdata_len > 0) then
				text = string.gsub(text, "]]>$", "")
			end

            if not string.find(text, "^%s*$") then
                local lVal = (top:value() or "") .. self:fromxmlstring(text)
                stack[#stack]:setValue(lVal)
            end
            if empty == "/" then -- empty element tag
                local lNode = newNode(label)
                self:ParseArgs(lNode, xarg)
                top:addChild(lNode)
            elseif c == "" then -- start tag
                local lNode = newNode(label)
                self:ParseArgs(lNode, xarg)
                table.insert(stack, lNode)
		top = lNode
            else -- end tag
                local toclose = table.remove(stack) -- remove top

                top = stack[#stack]
                if #stack < 1 then
                    error("xxml: nothing to close with " .. label)
                end
                if toclose:name() ~= label then
                    error("xxml: trying to close " .. toclose.name .. " with " .. label)
                end
                top:addChild(toclose)
            end
            i = j + 1
        end
        local text = string.sub(xmlText, i);
        if #stack > 1 then
            error("xxml: unclosed " .. stack[#stack]:name())
        end
        return top
    end

    function xxml:loadfile(path)

        local hFile, err = io.open(path, "r");

        if hFile and not err then
            local xmlText = hFile:read("*a"); -- read file content
            io.close(hFile);
            return self:load(xmlText), nil;
        else
            print(err)
            return nil
        end
    end



local _attrfuns={
_name=function(self) return self:name() end,
_value=function(self) return self:value() end,
innertext=function(self) return self:value() end,
tagname=function(self) return self:name() end}
function _getattr(tag,name)
	if(tag==nil) then
		return nil
	end
	local f=_attrfuns[string.lower(name)]
	if(f) then
		return f(tag)
	end
	return tag["@"..name]
end

function _current(tag,attrs)
	local attrs=attrs or {}
	for k,v in pairs(attrs) do
		local value=_getattr(tag,k)
		if(not(value) or value~=v) then
			return false
		end
	end
	return true
end



function _find(tags,attr)
	--1. 根据标签筛选
	local tag=tags[attr.tag]
	if(not(tag)) then
		return false
	end
	--2. 根据索引筛选
	if(attr.index) then
		tag=tag[attr.index]
	end
	if(not(tag)) then
		return false
	end
	local a=xtable.isarray(tag)
	local r={}
	if(a) then
		for i=1,#tag,1 do
			if(_current(tag[i],attr.attrs)) then
				table.insert(r,tag[i])
			end
		end

	else
		if(_current(tag,attr.attrs)) then
			table.insert(r,tag)
		end

	end
	--3. 从多个属性中筛选
	local cs=#r>1 and r or (#r>0 and r[1] or nil)
	return cs~=nil,cs
end


function newNode(name)
    local node = {}
    node.___value = nil
    node.___name = name
    node.___children = {}
    node.___props = {}


function node:get(path,name,i)
	local nodes=node:finds(path)
	if(xtable.isarray(nodes)) then
		return _getattr(nodes[i or 1],name)
	end
	return _getattr(nodes,name)
end

function node:finds(path)
	local pathtb=xpath.parse(path) or {}
	local c=self
	for i,v in pairs(pathtb) do
		local s,n=xxml._find(c,v)
		if(not(s)) then
			return node
		end
		c=n
	end
	return c or node
end
    function node:value() return self.___value end
    function node:setValue(val) self.___value = val end
    function node:name() return self.___name end
    function node:setName(name) self.___name = name end
    function node:children() return self.___children end
    function node:numChildren() return #self.___children end

    function node:addChild(child)
        if self[child:name()] ~= nil then
            if type(self[child:name()].name) == "function" then
                local tempTable = {}
                table.insert(tempTable, self[child:name()])
                self[child:name()] = tempTable
            end
            table.insert(self[child:name()], child)
        else
            self[child:name()] = child
        end
        table.insert(self.___children, child)
    end

    function node:properties() return self.___props end
    function node:numProperties() return #self.___props end
    function node:addProperty(name, value)
        local lName = "@" .. name
        if self[lName] ~= nil then
            if type(self[lName]) == "string" then
                local tempTable = {}
                table.insert(tempTable, self[lName])
                self[lName] = tempTable
            end
            table.insert(self[lName], value)
        else
            self[lName] = value
        end
        table.insert(self.___props, { name = name, value = self[name] })
    end
	return node
end

--local nodes=xxml:loadfile("c:\\Cache.Config")
--print(nodes:get("/caches/cache[name=comm]/item[name]","key",1))


