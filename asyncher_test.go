package goasync

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	a := New(
		func() (interface{}, error) {
			buf := new(bytes.Buffer)

			buf.WriteString("main\n\r")
			return buf, nil
		},
		func(i interface{}, e error) (interface{}, error) {
			buf := i.(*bytes.Buffer)
			buf.WriteString("finally\n\r")
			return buf, nil
		},
		nil, // no panic function
	)

	i, _ := a.Process()
	buf := i.(*bytes.Buffer)

	assert.Equal(t, buf.String(), "main\n\rfinally\n\r")
}

func TestAsync(t *testing.T) {
	a := New(
		func() (interface{}, error) {
			buf := new(bytes.Buffer)

			buf.WriteString("main\n\r")
			return buf, nil
		},
		func(i interface{}, e error) (interface{}, error) {
			buf := i.(*bytes.Buffer)
			buf.WriteString("finally\n\r")
			return buf, nil
		},
		nil, // no panic function
	)

	w := a.ProcessAsync()
	i, _ := w.Wait()
	buf := i.(*bytes.Buffer)

	assert.Equal(t, buf.String(), "main\n\rfinally\n\r")
}

func TestSyncClosure(t *testing.T) {
	buf := new(bytes.Buffer)

	a := New(
		func() (interface{}, error) {
			buf.WriteString("main\n\r")
			return nil, nil
		},
		func(interface{}, error) (interface{}, error) {
			buf.WriteString("finally\n\r")
			return nil, nil
		},
		nil, // no panic function
	)

	a.Process()

	assert.Equal(t, buf.String(), "main\n\rfinally\n\r")
}

func TestAsyncClosure(t *testing.T) {
	buf := new(bytes.Buffer)

	a := New(
		func() (interface{}, error) {
			buf.WriteString("main\n\r")
			return nil, nil
		},
		func(interface{}, error) (interface{}, error) {
			buf.WriteString("finally\n\r")
			return nil, nil
		},
		nil, // no panic function
	)

	waiter := a.ProcessAsync()

	waiter.Wait()

	assert.Equal(t, buf.String(), "main\n\rfinally\n\r")
}
