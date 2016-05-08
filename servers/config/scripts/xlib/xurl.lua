
module ("xurl", package.seeall)

function encode(s)  
    local str=s
  if (str) then  
    str = string.gsub (str, "\n", "\r\n")  
    str = string.gsub (str, "([^%w ])",  
        function (c) return string.format ("%%%02X", string.byte(c)) end)  
    str = string.gsub (str, " ", "+")  
  end  
  return str      
end  
  
function decode(s)  
 local str = string.gsub (s, "+", " ")  
  str = string.gsub (str, "%%(%x%x)",  
      function(h) return string.char(tonumber(h,16)) end)  
  str = string.gsub (str, "\r\n", "\n")  
  return str  
end


