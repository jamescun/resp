package resp

import (
	"bufio"
	"bytes"
	"io"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// ReadLine reads a line from the buffer and returns the lines type and value.
// The returned values are slices of the buffer and are only valid until the next read,
// you should take care when using its return value or use one of the helper methods
// for specific types.
func (r *Reader) ReadLine() (t byte, v []byte, err error) {
	b, err := r.r.ReadSlice('\n')
	if err != nil {
		return
	} else if len(b) < 4 {
		err = ErrProtocol
		return
	}

	t = b[0]
	v = b[1 : len(b)-2]

	return
}

// NextType peeks ahead in the buffer to the next line and returns the type,
// it will block if no line is available.
func (r *Reader) NextType() (byte, error) {
	b, err := r.r.Peek(1)
	if err != nil {
		return 0x00, err
	}

	return b[0], nil
}

func (r *Reader) ReadString() (string, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return "", err
	} else if t != String {
		return "", ErrUnexpectedType
	}

	// optimisation to avoid allocations for frequent PING/PONG messages
	if len(v) == 4 {
		if bytes.Equal(v, []byte(PING)) {
			return PING, nil
		} else if bytes.Equal(v, []byte(PONG)) {
			return PONG, nil
		}
	}

	return string(v), nil
}

func (r *Reader) ReadError() (string, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return "", err
	} else if t != Error {
		return "", ErrUnexpectedType
	}

	return string(v), nil
}

func (r *Reader) ReadBulkString() ([]byte, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return nil, err
	} else if t != BulkString {
		return nil, ErrUnexpectedType
	}

	n, err := parseInt(v)
	if err != nil {
		return nil, err
	}

	b := make([]byte, n+2)
	_, err = io.ReadFull(r.r, b)
	if err != nil {
		return nil, err
	}

	return b[:len(b)-2], nil
}

func (r *Reader) BulkStringReader() (int, io.Reader, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return 0, nil, err
	} else if t != BulkString {
		return 0, nil, ErrUnexpectedType
	}

	n, err := parseInt(v)
	if err != nil {
		return 0, nil, err
	}
	n += 2

	return n, io.LimitReader(r.r, int64(n)), nil
}

func (r *Reader) ReadInteger() (int, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return 0, err
	} else if t != Integer {
		return 0, ErrUnexpectedType
	}

	return parseInt(v)
}

func (r *Reader) ReadArray() (int, error) {
	t, v, err := r.ReadLine()
	if err != nil {
		return 0, err
	} else if t != Array {
		return 0, ErrUnexpectedType
	}

	return parseInt(v)
}
