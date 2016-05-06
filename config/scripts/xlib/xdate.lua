require 'xnumber'
require 'xtimer'
require 'xtable'

module ("xdate", package.seeall)
local HOURPERDAY  = 24
local MINPERHOUR  = 60
local MINPERDAY	  = 1440  -- 24*60
local SECPERMIN   = 60
local SECPERHOUR  = 3600  -- 60*60
local SECPERDAY   = 86400 -- 24*60*60
local TICKSPERSEC = 1000000
local TICKSPERDAY = 86400000000
local TICKSPERHOUR = 3600000000
local TICKSPERMIN = 60000000
local DAYNUM_MAX =  365242500 -- Sat Jan 01 1000000 00:00:00
local DAYNUM_MIN = -365242500 -- Mon Jan 01 1000000 BCE 00:00:00
local DAYNUM_DEF =  0 -- Mon Jan 01 0001 00:00:00
local floor = math.floor
local ceil  = math.ceil
local abs   = math.abs
local sub  = string.sub
local gsub = string.gsub
local gmatch = string.gmatch or string.gfind
local find = string.find
local function mod (n,d) return n - d*floor(n/d) end
local fmt  = string.format
local fmtstr  = "%Y-%M-%d %H:%m:%S"



local formatarg={f={yyyy="%%Y"},s={yy="%%y",MM="%%M",dd="%%d",hh="%%H",HH="%%H",mm="%%m",mi="%%m",ss="%%S",wd="%%a"}}
xdate._formatReplace=function(f)
	local r=f
	for k,v in pairs(formatarg.f) do
		r=string.gsub(r,k,v)
	end
	for k,v in pairs(formatarg.s) do
		r=string.gsub(r,k,v)
	end
	return r
end


local sl_weekdays = {
		[0]="周日",[1]="周一",[2]="周二",[3]="周三",[4]="周四",[5]="周五",[6]="周六"
	}
local sl_meridian = {[-1]="AM", [1]="PM"}
local sl_months = {
	[00]="1月", [01]="2月", [02]="3月",
	[03]="4月",   [04]="5月",      [05]="6月",
	[06]="7月",    [07]="8月",   [08]="9月",
	[09]="10月", [10]="11月", [11]="12月"
}
local dayfromyear=function(y)
	return 365*y + floor(y/4) - floor(y/100) + floor(y/400)
end

local makedaynum=function(y, m, d)
	local mm =mod(mod(m,12) + 10, 12)
	return dayfromyear(y + floor(m/12) - floor(mm/10)) + floor((mm*306 + 5)/10) + d - 307
end

local weekday=function(dn)
	return mod(dn + 1, 7)
end
local date_error_arg=function() return error("invalid argument(s)",0) end
local breakdaynum=function(g)
	local g = g + 306
	local y = floor((10000*g + 14780)/3652425)
	local d = g - dayfromyear(y)
	if d < 0 then y = y - 1; d = g - dayfromyear(y) end
	local mi = floor((100*d + 52)/3060)
	return (floor((mi + 2)/12) + y),mod(mi + 2,12), (d - floor((mi*306 + 5)/10) + 1)
end
local makedayfrc=function(h,r,s,t)
		return ((h*60 + r)*60 + s)*TICKSPERSEC + t
end
local function fix(n) n = tonumber(n) return n and ((n > 0 and floor or ceil)(n)) end
local function getmontharg(v)
		local m = tonumber(v)
		return (m and fix(m - 1)) or inlist(tostring(v) or "", sl_months, 2)
	end
local function breakdayfrc(df)
	return
		mod(floor(df/TICKSPERHOUR),HOURPERDAY),
		mod(floor(df/TICKSPERMIN ),MINPERHOUR),
		mod(floor(df/TICKSPERSEC ),SECPERMIN),
		mod(df,TICKSPERSEC)
end


function xdate:normalize()
	local dn, df = fix(self.daynum), self.dayfrc
	self.daynum, self.dayfrc = dn + floor(df/TICKSPERDAY),mod(df, TICKSPERDAY)
	return (dn >= DAYNUM_MIN and dn <= DAYNUM_MAX) and self or error("date beyond imposed limits:"..self)
end

function xdate:addyears(y, m, d)
	local cy, cm, cd = breakdaynum(self.daynum)
	if y then y = fix(tonumber(y))else y = 0 end
	if m then m = fix(tonumber(m))else m = 0 end
	if d then d = fix(tonumber(d))else d = 0 end
	if y and m and d then
		self.daynum  = makedaynum(cy+y, cm+m, cd+d)
		return self:normalize()
	else
		return date_error_arg()
	end
end
function xdate:addmonths(m, d)
	return self:addyears(nil, m, d)
end

local function xdate_adddayfrc(self,n,pt,pd)
	n = tonumber(n)
	if n then
		local x = floor(n/pd);
		self.daynum = self.daynum + x;
		self.dayfrc = self.dayfrc + (n-x*pd)*pt;
		return self:normalize()
	else
		return date_error_arg()
	end
end
function xdate:adddays(n)	return xdate_adddayfrc(self,n,TICKSPERDAY,1) end
function xdate:addhours(n)	return xdate_adddayfrc(self,n,TICKSPERHOUR,HOURPERDAY) end
function xdate:addminutes(n)	return xdate_adddayfrc(self,n,TICKSPERMIN,MINPERDAY)  end
function xdate:addseconds(n)	return xdate_adddayfrc(self,n,TICKSPERSEC,SECPERDAY)  end
function xdate:addticks(n)	return xdate_adddayfrc(self,n,1,TICKSPERDAY) end
function xdate:getdate()	local y, m, d = breakdaynum(self.daynum) return y, m+1, d end
function xdate:gettime()	return breakdayfrc(self.dayfrc) end
function xdate:getclockhour() local h = self:gethours() return h>12 and mod(h,12) or (h==0 and 12 or h) end
function xdate:getyearday() return yearday(self.daynum) + 1 end
function xdate:getweekday() return weekday(self.daynum) + 1 end   -- in lua weekday is sunday = 1, monday = 2 ...
function xdate:getyear()	 local r,_,_ = breakdaynum(self.daynum)	return r end
function xdate:getmonth() local _,r,_ = breakdaynum(self.daynum)	return r+1 end-- in lua month is 1 base
function xdate:getday()	 local _,_,r = breakdaynum(self.daynum)	return r end
function xdate:gethours()	return mod(floor(self.dayfrc/TICKSPERHOUR),HOURPERDAY) end
function xdate:getminutes()	return mod(floor(self.dayfrc/TICKSPERMIN), MINPERHOUR) end
function xdate:getseconds()	return mod(floor(self.dayfrc/TICKSPERSEC ),SECPERMIN)  end
function xdate:getfracsec()	return mod(floor(self.dayfrc/TICKSPERSEC ),SECPERMIN)+(mod(self.dayfrc,TICKSPERSEC)/TICKSPERSEC) end
function xdate:getticks(u)	local x =mod(self.dayfrc,TICKSPERSEC) return u and ((x*u)/TICKSPERSEC) or x  end

 xdate.tpf={
 d=function(self,v) self:adddays(v)end,
 m=function(self,v) self:addmonths(v)end,
 y=function(self,v) self:addyears(v)end,
 h=function(self,v) self:addhours(v)end}

function xdate:getweeknumber(wdb)
	local wd, yd = weekday(self.daynum), yearday(self.daynum)
	if wdb then
		wdb = tonumber(wdb)
		if wdb then
			wd =mod(wd-(wdb-1),7)-- shift the week day base
		else
			return date_error_arg()
		end
	end
	return (yd < wd and 0) or (floor(yd/7) + ((mod(yd, 7)>=wd) and 1 or 0))
end
xdate.tvspec = {
		-- Full weekday name (Sunday)
		['%a']=function(self) return sl_weekdays[weekday(self.daynum)] end,
		-- Full month name (December)
		['%b']=function(self) return sl_months[self:getmonth() - 1] end,
		-- The day of the month as a number (range 1 - 31)
		['%d']=function(self) return fmt("%.2d", self:getday())  end,
		-- hour of the 24-hour day, from 00 (06)
		['%H']=function(self) return fmt("%.2d", self:gethours()) end,
		-- The  hour as a number using a 12-hour clock (01 - 12)
		['%I']=function(self) return fmt("%.2d", self:getclockhour()) end,
		-- The day of the year as a number (001 - 366)
		['%j']=function(self) return fmt("%.3d", self:getyearday())  end,
		-- Month of the year, from 01 to 12
		['%M']=function(self) return fmt("%.2d", self:getmonth())  end,
		-- Minutes after the hour 55
		['%m']=function(self) return fmt("%.2d", self:getminutes())end,
		-- AM/PM indicator (AM)
		['%p']=function(self) return sl_meridian[self:gethours() > 11 and 1 or -1] end, --AM/PM indicator (AM)
		-- The second as a number (59, 20 , 01)
		['%S']=function(self) return fmt("%.2d", self:getseconds())  end,
		-- Sunday week of the year, from 00 (48)
		['%U']=function(self) return fmt("%.2d", self:getweeknumber()) end,
		-- The day of the week as a decimal, Sunday being 0
		['%w']=function(self) return self:getweekday() - 1 end,
		-- Monday week of the year, from 00 (48)
		['%W']=function(self) return fmt("%.2d", self:getweeknumber(2)) end,
		-- The year as a number without a century (range 00 to 99)
		['%y']=function(self) return fmt("%.2d", mod(self:getyear() ,100)) end,
		-- Year with century (2000, 1914, 0325, 0001)
		['%Y']=function(self) return fmt("%.4d", self:getyear()) end

	}
	function xdate:fmt0(str) return (gsub(str, "%%[%a%%\b\f]", function(x) local f = xdate.tvspec[x];return (f and f(self)) or x end)) end
	function xdate:fmt(str)
		str = str or self.fmtstr or fmtstr
		return self:fmt0((gmatch(str, "${%w+}")) and (gsub(str, "${%w+}", function(x)local f=tvspec[x];return (f and f(self)) or x end)) or str)
	end

	function xdate.__lt(a,b)	return (a.daynum == b.daynum) and (a.dayfrc < b.dayfrc) or (a.daynum < b.daynum)	end
	function xdate.__le(a, b)return (a.daynum == b.daynum) and (a.dayfrc <= b.dayfrc) or (a.daynum <= b.daynum)	end
	function xdate.__eq(a, b)return (a.daynum == b.daynum) and (a.dayfrc == b.dayfrc) end
	function xdate.__sub(a,b)
		local d1, d2 = date_getxdate(a), date_getxdate(b)
		local d0 = d1 and d2 and date_new(d1.daynum - d2.daynum, d1.dayfrc - d2.dayfrc)
		return d0 and d0:normalize()
	end
	function xdate.__add(a,b)
		local d1, d2 = date_getxdate(a), date_getxdate(b)
		local d0 = d1 and d2 and date_new(d1.daynum + d2.daynum, d1.dayfrc + d2.dayfrc)
		return d0 and d0:normalize()
	end
	function xdate.__concat(a, b) return tostring(a) .. tostring(b) end
	function xdate:__tostring() return self:fmt() end




function xdate:setyear(y, m, d)
	local cy, cm, cd = breakdaynum(self.daynum)
	if y then cy = fix(tonumber(y))end
	if m then cm = getmontharg(m)  end
	if d then cd = fix(tonumber(d))end
	if cy and cm and cd then
		self.daynum  = makedaynum(cy, cm, cd)
		return self:normalize()
	else
		return date_error_arg()
	end
end

function xdate:sethours(h, m, s, t)
	local ch,cm,cs,ck = breakdayfrc(self.dayfrc)
	ch, cm, cs, ck = tonumber(h or ch), tonumber(m or cm), tonumber(s or cs), tonumber(t or ck)
	if ch and cm and cs and ck then
		self.dayfrc = makedayfrc(ch, cm, cs, ck)
		return self:normalize()
	else
		return date_error_arg()
	end
end
function xdate:format(f)
	local fm=xdate._formatReplace(f or fmtstr)
	return self:fmt(fm)
end

function xdate:now(f)
	local o={}
	o.time=os.time()
	o.year=os.date("%Y",o.time)
	o.month=os.date("%m",o.time)
	o.day=os.date("%d",o.time)
	o.hour=os.date("%H",o.time)
	o.minute=os.date("%M",o.time)
	o.second=os.date("%S",o.time)
    setmetatable(o,self)
    self.__index = self
	self.fmtstr=xdate._formatReplace(f or fmtstr)
	o.daynum=makedaynum(o.year, o.month - 1, o.day)
	o.dayfrc= makedayfrc(o.hour, o.minute,o.second, 0)
	return o
end
function xdate:new(d,fm)
	local f=fm or "yyyy-MM-dd hh:mm:ss"
	local o={}
	o.year=string.find(f,"yyyy")~=nil and string.sub(d,string.find(f,"yyyy")) or 0
	o.month=string.find(f,"MM")~=nil and string.sub(d,string.find(f,"MM")) or 0
	o.day=string.find(f,"dd")~=nil and string.sub(d,string.find(f,"dd")) or 0
	o.hour=string.find(f,"hh")~=nil and string.sub(d,string.find(f,"hh")) or 0
	o.minute=string.find(f,"mm")~=nil and string.sub(d,string.find(f,"mm")) or 0
	o.second=string.find(f,"ss")~=nil and string.sub(d,string.find(f,"ss")) or 0
	o.year=o.year==nil and (string.find(f,"yy")~=nil and string.sub(d,string.find(f,"yy")) or 0) or o.year
	o.hour=o.hour==nil and (string.find(f,"HH")~=nil and string.sub(d,string.find(f,"HH")) or 0) or o.hour
	o.minute=o.minute==nil and (string.find(f,"mi")~=nil and string.sub(d,string.find(f,"mi")) or 0) or o.minute
	setmetatable(o,self)
    self.__index = self
	o.daynum=makedaynum(o.year, o.month - 1, o.day)
	o.dayfrc= makedayfrc(o.hour, o.minute,o.second, 0)
	return o
end

xdate.get_range=function(d,t,c)
	local r={}
	local tpf=xdate.tpf[t or "d"]
	assert(tpf~=nil,"传入的类型不合法")
	local v=d>0 and 1 or -1
	local len=xnumber.parse(math.abs(d),0)
	local _date=c or xdate:now()
	for i=1,len,1 do
		local cdate=xtable.clone(_date)
		tpf(cdate,v*(i-1))
		table.insert(r,cdate)
	end
	return r
end

xdate.get_month_days=function(m)
	local r={}
	local _date=m
	if(not(_date)) then
		_=xdate:now()
		_date=xdate:new(string.format("%s-%s-%s",_.year,_.month,"01"),"yyyy-MM-dd")
	end
	local end_date=xtable.clone(_date):addmonths(1)

	 while(_date<end_date) do
		table.insert(r,_date)
		_date=xtable.clone(_date):adddays(1)
	end
	return r
end

