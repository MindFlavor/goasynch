package async

import (
	"bytes"
	"fmt"
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

	assert.Equal(t, "main\n\rfinally\n\r", buf.String())
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

	assert.Equal(t, "main\n\rfinally\n\r", buf.String())
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

	assert.Equal(t, "main\n\rfinally\n\r", buf.String())
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

	assert.Equal(t, "main\n\rfinally\n\r", buf.String())
}

func TestPanicClosure(t *testing.T) {
	buf := new(bytes.Buffer)

	a := New(
		func() (interface{}, error) {
			buf.WriteString("main\n\r")
			panic("!!!")

		},
		func(interface{}, error) (interface{}, error) {
			buf.WriteString("finally\n\r")
			return nil, nil
		},
		func(e interface{}) {
			buf.WriteString(fmt.Sprintf("recover%s\n\r", e))
		},
	)

	a.Process()

	assert.Equal(t, "main\n\rrecover!!!\n\rfinally\n\r", buf.String())
}

func TestPanicClosureFinallyOnlyOnce(t *testing.T) {
	buf := new(bytes.Buffer)

	a := New(
		func() (interface{}, error) {
			buf.WriteString("main\n\r")
			panic("!!!")

		},
		func(interface{}, error) (interface{}, error) {
			buf.WriteString("finally\n\r")
			panic("$$$")
		},
		func(e interface{}) {
			buf.WriteString(fmt.Sprintf("recover%s\n\r", e))
		},
	)

	a.Process()

	assert.Equal(t, "main\n\rrecover!!!\n\rfinally\n\rrecover$$$\n\r", buf.String())
}
