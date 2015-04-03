package wolf

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBodyWrapper(t *testing.T) {
	var b io.ReadCloser = &bodyWrapper{}

	assert.Panics(t, func() {
		b.Read(nil)
	})
	assert.Panics(t, func() {
		b.Close()
	})
}
