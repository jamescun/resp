package resp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	n, err := parseInt([]byte("100"))
	if assert.NoError(t, err) {
		assert.Equal(t, 100, n)
	}

	n, err = parseInt([]byte("-200"))
	if assert.NoError(t, err) {
		assert.Equal(t, -200, n)
	}

	_, err = parseInt([]byte(""))
	assert.Equal(t, ErrShort, err)

	_, err = parseInt([]byte("-"))
	assert.Equal(t, ErrShort, err)

	_, err = parseInt([]byte("foo"))
	assert.Equal(t, ErrProtocol, err)

	_, err = parseInt([]byte("100\r\n"))
	assert.Equal(t, ErrProtocol, err)
}
