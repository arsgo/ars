package script

import lua "github.com/yuin/gopher-lua"

//Post(url string, params string, args ...string) (content string, status int, err error)
func getStringParams(ls *lua.LState, start int) (params []string) {
	c := ls.GetTop()
	params = make([]string, 0, c)
	for i := start; i <= c; i++ {
		t := ls.Get(i).Type().String()
		if t == "userdata" {
			ls.RaiseError("invalid string of function arguments (string expected, got userdata)")
		} else {
			params = append(params, ls.Get(i).String())
		}
	}
	return
}
func getMapParams(tb *lua.LTable) map[string]string {
	data := make(map[string]string)
	tb.ForEach(func(key lua.LValue, value lua.LValue) {
		if value != nil {
			data[key.String()] = value.String()
		}
	})
	return data
}
