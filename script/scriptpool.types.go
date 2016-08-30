package script

import "github.com/arsgo/lib4go/script"

func (s *ScriptPool) bindTypes() (r []script.LuaTypesBinder) {
	r = append(r, s.getHttpTypeBinder())
	r = append(r, s.getMemcachedBinder())
	r = append(r, s.getWeixinTypeBinder())
	r = append(r, s.getElasticTypeBinder())
	r = append(r, s.getinfluxTypeBinder())
	r = append(r, s.getMQTypeBinder())
	return

}
