package script

import "github.com/yuin/gopher-lua"

func (s *ScriptPool) moduleMQProducerSend(ls *lua.LState) int {
	config := ls.CheckString(1)
	queue := ls.CheckString(2)
	content := ls.CheckTable(3)
	timeout := 0
	if ls.GetTop() == 4 {
		timeout = ls.CheckInt(4)
	}
	pdc, err := s.NewMQProducer(config)
	if err != nil {
		return pushValues(ls, err)
	}
	tb, err := luaTable2Json(content)
	if err != nil {
		return pushValues(ls, err.Error())
	}
	err = pdc.Send(queue, tb, timeout)
	return pushValues(ls, err)
}
