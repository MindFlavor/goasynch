package goasync

// The Waiter interface allows you to
// synchronize the async call and
// retrieve the results.
type Waiter interface {
	Wait() (interface{}, error)
}

type wait struct {
	result  interface{}
	err     error
	cFinish chan (bool)
}

func newWaiter() *wait {
	return &wait{
		result:  nil,
		err:     nil,
		cFinish: make(chan (bool), 1),
	}
}

func (a *wait) Wait() (interface{}, error) {
	<-a.cFinish

	return a.result, a.err
}
