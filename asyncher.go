// Package goasync incapsulates a
// long running goroutine with sync
// at the end
package goasync

// Asyncer interface exposes the methods
// you can call start the processing.
// The package will call these function in order:
// mainFnc and then finallyFnc. If mainFnc panics the flow will be
// mainFnc -> panicFnc -> finallyFnc.
// If finallyFnc panics the flow will be:
// mainFnc -> finallyFnc -> panicFnc.
// If finallyFnc panics it will not be called again after panicFnc.
type Asyncer interface {
	Process() (interface{}, error)
	ProcessAsync() Waiter
}

type selfdisposer struct {
	finallyFnc func(interface{}, error) (interface{}, error)
	mainFnc    func() (interface{}, error)
	panicFnc   func(e interface{})
}

// New creates a new idle Asyncher.
// To start it call its methods.
func New(mainFnc func() (interface{}, error), finallyFnc func(interface{}, error) (interface{}, error), panicFnc func(e interface{})) Asyncer {
	return &selfdisposer{
		finallyFnc: finallyFnc,
		mainFnc:    mainFnc,
		panicFnc:   panicFnc,
	}
}

func (s *selfdisposer) Process() (interface{}, error) {
	fFinal := false
	if s.panicFnc != nil {
		defer func() {
			if e := recover(); e != nil {
				s.panicFnc(e)
			}

			if !fFinal {
				defer func() {
					if e := recover(); e != nil {
						s.panicFnc(e)
					}
				}()
				s.finallyFnc(nil, nil)
			}
		}()
	}

	ret, err := s.mainFnc()

	fFinal = true
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
