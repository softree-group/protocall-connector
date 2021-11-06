package bus

import "protocall/domain/services"

type Subscriber struct {
	C     chan interface{}
	clear func()
}

func (s *Subscriber) Cancel() {
	close(s.C)
	s.C = nil

	s.clear()
}

func (s Subscriber) Channel() chan interface{} {
	return s.C
}

var _ services.Subscriber = &Subscriber{}
