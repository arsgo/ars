package rpcservice


type ServiceCallBack interface {
	Request(string, string) (string, error)
	Send(string, string, []byte) (string, error)
    Get(string, string) ([]byte, error)
}

type serviceHandler struct {
	services map[string]*spSnap
	callback ServiceCallBack
}

func (s *serviceHandler) Request(name string, input string) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &spSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return s.callback.Request(name, input)
}

func (s *serviceHandler) Send(name string, input string, data []byte) (r string, err error) {
	if _, ok := s.services[name]; !ok {
		s.services[name] = &spSnap{}
	}
	s.services[name].total++
	s.services[name].current++
	defer func() {
		s.services[name].current--
		if err == nil {
			s.services[name].success++
		} else {
			s.services[name].failed++
		}
	}()
	return s.callback.Send(name, input,data)
}

func NewServiceHandler(callback ServiceCallBack) *serviceHandler {
	return &serviceHandler{services: make(map[string]*spSnap), callback: callback}
}