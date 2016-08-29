package script

import lua "github.com/yuin/gopher-lua"

//report.success 上报成功结果
func (s *ScriptPool) moduleReportSuccess(ls *lua.LState) int {
	name := ls.CheckString(1)
	taskType := ls.GetGlobal("__task_type__").String()
	taskName := ls.GetGlobal("__task_name__").String()
	if c, ok := s.collectors[taskType]; ok {
		c.Customer(taskName).Success(name)
		return pushValues(ls, "")
	}
	return pushValues(ls, "未找到report collector")
}

//report.error 上报错误
func (s *ScriptPool) moduleReportError(ls *lua.LState) int {
	name := ls.CheckString(1)
	taskType := ls.GetGlobal("__task_type__").String()
	taskName := ls.GetGlobal("__task_name__").String()
	if c, ok := s.collectors[taskType]; ok {
		c.Customer(taskName).Error(name)
		return pushValues(ls, "")
	}
	return pushValues(ls, "未找到report collector")
}

//report.failed 上报失败
func (s *ScriptPool) moduleReportFaild(ls *lua.LState) int {
	name := ls.CheckString(1)
	taskType := ls.GetGlobal("__task_type__").String()
	taskName := ls.GetGlobal("__task_name__").String()
	if c, ok := s.collectors[taskType]; ok {
		c.Customer(taskName).Failed(name)
		return pushValues(ls, "")
	}
	return pushValues(ls, "未找到report collector")
}

//report.Juge 上报失败
func (s *ScriptPool) moduleReportJuge(ls *lua.LState) int {
	value := ls.CheckBool(1)
	name := ls.CheckString(2)
	taskType := ls.GetGlobal("__task_type__").String()
	taskName := ls.GetGlobal("__task_name__").String()
	if c, ok := s.collectors[taskType]; ok {
		c.Customer(taskName).Juge(value, name)
		return pushValues(ls, "")
	}
	return pushValues(ls, "未找到report collector")
}
