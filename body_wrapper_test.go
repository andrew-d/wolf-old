package wolf

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyReadCloser struct {
	readCalled  bool
	closeCalled bool
}

func (d *dummyReadCloser) Read(buf []byte) (int, error) {
	d.readCalled = true
	return 0, nil
}

func (d *dummyReadCloser) Close() error {
	d.closeCalled = true
	return nil
}

func TestBodyWrapper(t *testing.T) {
	u := &dummyReadCloser{}
	var b io.ReadCloser = &bodyWrapper{underlying: u}

	b.Read(nil)
	assert.True(t, u.readCalled)

	b.Close()
	assert.True(t, u.closeCalled)
}
