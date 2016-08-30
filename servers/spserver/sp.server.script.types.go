package main

import "github.com/arsgo/lib4go/script"

func (s *SPServer) bindTypes() (r []script.LuaTypesBinder) {
	r = append(r, s.getdbTypeBinder()...)
	return

}
