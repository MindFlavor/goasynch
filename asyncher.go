// Package selfdisposer incapsulates a
// long running goroutine with disposing (closing)
// at the end
package goasync

type Asyncer interface {
	Process() (interface{}, error)
	ProcessAsync() Waiter
}

type selfdisposer struct {
	finallyFnc func(interface{}, error) (interface{}, error)
	mainFnc    func() (interface{}, error)
	panicFnc   func()
}

func New(mainFnc func() (interface{}, error), finallyFnc func(interface{}, error) (interface{}, error), panicFnc func()) Asyncer {
	return &selfdisposer{
		finallyFnc: finallyFnc,
		mainFnc:    mainFnc,
		panicFnc:   panicFnc,
	}
}

func (s *selfdisposer) Process() (interface{}, error) {
	if s.panicFnc != nil {
		defer s.panicFnc()
	}

	ret, err := s.mainFnc()

	if s.finallyFnc != nil {
		return s.finallyFnc(ret, err)
	}

	return ret, err
}

func (s *selfdisposer) ProcessAsync() Waiter {
	a := newWaiter()

	go func() {
		a.result, a.err = s.Process()
		a.cFinish <- true
	}()

	return a
}
