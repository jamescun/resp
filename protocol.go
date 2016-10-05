package resp

import (
	"errors"
)

var (
	ErrShort          = errors.New("ERR Short")
	ErrProtocol       = errors.New("ERR Protocol Error")
	ErrUnexpectedType = errors.New("ERR Unexpected Type")
)

const (
	String     byte = '+'
	Error      byte = '-'
	Integer    byte = ':'
	BulkString byte = '$'
	Array      byte = '*'

	OK   string = "OK"
	PING string = "PING"
	PONG string = "PONG"
	CRLF string = "\r\n"
)

func parseInt(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, ErrShort
	}

	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, ErrShort
		}
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, ErrProtocol
		}
		n += int(b - '0')
	}

	if negate {
		n = -n
	}

	return n, nil
}
