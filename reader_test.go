package resp

import (
	"bufio"
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderNew(t *testing.T) {
	r := NewReader(bytes.NewReader([]byte("")))
	if assert.NotNil(t, r) {
		assert.NotNil(t, r.r)
	}
}

func TestReaderReadLine(t *testing.T) {
	r := newReader(":100\r\n:-200\r\n:\r\n")

	m, v, err := r.ReadLine()
	if assert.NoError(t, err) {
		assert.Equal(t, Integer, m)
		assert.Equal(t, []byte("100"), v)
	}

	m, v, err = r.ReadLine()
	if assert.NoError(t, err) {
		assert.Equal(t, Integer, m)
		assert.Equal(t, []byte("-200"), v)
	}

	_, _, err = r.ReadLine()
	assert.Equal(t, ErrProtocol, err)

	_, _, err = r.ReadLine()
	assert.Equal(t, io.EOF, err)
}

func TestReaderNextType(t *testing.T) {
	r := newReader(":100\r\n")

	m, err := r.NextType()
	if assert.NoError(t, err) {
		assert.Equal(t, Integer, m)
	}

	r = newReader("")

	_, err = r.NextType()
	assert.Equal(t, io.EOF, err)
}

func TestReaderReadString(t *testing.T) {
	r := newReader("+OK\r\n+Foo\r\n+PING\r\n+PONG\r\n-ERR\r\n+\r\n")

	s, err := r.ReadString()
	if assert.NoError(t, err) {
		assert.Equal(t, "OK", s)
	}

	s, err = r.ReadString()
	if assert.NoError(t, err) {
		assert.Equal(t, "Foo", s)
	}

	s, err = r.ReadString()
	if assert.NoError(t, err) {
		assert.Equal(t, "PING", s)
	}

	s, err = r.ReadString()
	if assert.NoError(t, err) {
		assert.Equal(t, "PONG", s)
	}

	_, err = r.ReadString()
	assert.Equal(t, ErrUnexpectedType, err)

	_, err = r.ReadString()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadString()
	assert.Equal(t, io.EOF, err)
}

func TestReaderReadError(t *testing.T) {
	r := newReader("-ERR\r\n-Foo\r\n+OK\r\n-\r\n")

	s, err := r.ReadError()
	if assert.NoError(t, err) {
		assert.Equal(t, "ERR", s)
	}

	s, err = r.ReadError()
	if assert.NoError(t, err) {
		assert.Equal(t, "Foo", s)
	}

	_, err = r.ReadError()
	assert.Equal(t, ErrUnexpectedType, err)

	_, err = r.ReadError()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadError()
	assert.Equal(t, io.EOF, err)
}

func TestReaderReadBulkString(t *testing.T) {
	r := newReader("$3\r\nfoo\r\n$3\r\nbar\r\n+foo\r\n$foo\r\n$\r\n$4\r\nbaz\r\n")

	b, err := r.ReadBulkString()
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("foo"), b)
	}

	b, err = r.ReadBulkString()
	if assert.NoError(t, err) {
		assert.Equal(t, []byte("bar"), b)
	}

	_, err = r.ReadBulkString()
	assert.Equal(t, ErrUnexpectedType, err)

	_, err = r.ReadBulkString()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadBulkString()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadBulkString()
	assert.Equal(t, io.ErrUnexpectedEOF, err)

	_, err = r.ReadBulkString()
	assert.Equal(t, io.EOF, err)
}

func TestReaderBulkStringReader(t *testing.T) {
	r := newReader("$3\r\nfoo\r\n$3\r\nbar\r\n+foo\r\n$foo\r\n$\r\n")

	n, rd, err := r.BulkStringReader()
	if assert.NoError(t, err) {
		assert.Equal(t, 5, n)

		b := make([]byte, 5)
		_, err := io.ReadFull(rd, b)
		if assert.NoError(t, err) {
			assert.Equal(t, []byte("foo\r\n"), b)
		}
	}

	n, rd, err = r.BulkStringReader()
	if assert.NoError(t, err) {
		assert.Equal(t, 5, n)

		b := make([]byte, 5)
		_, err := io.ReadFull(rd, b)
		if assert.NoError(t, err) {
			assert.Equal(t, []byte("bar\r\n"), b)
		}
	}

	_, _, err = r.BulkStringReader()
	assert.Equal(t, ErrUnexpectedType, err)

	_, _, err = r.BulkStringReader()
	assert.Equal(t, ErrProtocol, err)

	_, _, err = r.BulkStringReader()
	assert.Equal(t, ErrProtocol, err)

	_, _, err = r.BulkStringReader()
	assert.Equal(t, io.EOF, err)
}

func TestReaderReadInteger(t *testing.T) {
	r := newReader(":100\r\n:-200\r\n-200\r\n:\r\n")

	n, err := r.ReadInteger()
	if assert.NoError(t, err) {
		assert.Equal(t, 100, n)
	}

	n, err = r.ReadInteger()
	if assert.NoError(t, err) {
		assert.Equal(t, -200, n)
	}

	_, err = r.ReadInteger()
	assert.Equal(t, ErrUnexpectedType, err)

	_, err = r.ReadInteger()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadInteger()
	assert.Equal(t, io.EOF, err)
}

func TestReaderReadArray(t *testing.T) {
	r := newReader("*2\r\n*3\r\n+OK\r\n*\r\n")

	n, err := r.ReadArray()
	if assert.NoError(t, err) {
		assert.Equal(t, 2, n)
	}

	n, err = r.ReadArray()
	if assert.NoError(t, err) {
		assert.Equal(t, 3, n)
	}

	_, err = r.ReadArray()
	assert.Equal(t, ErrUnexpectedType, err)

	_, err = r.ReadArray()
	assert.Equal(t, ErrProtocol, err)

	_, err = r.ReadArray()
	assert.Equal(t, io.EOF, err)
}

func newReader(s string) *Reader {
	return &Reader{r: bufio.NewReader(bytes.NewBufferString(s))}
}
