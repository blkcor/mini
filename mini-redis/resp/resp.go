package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// RESP protocol symbol
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Value is a RESP value
type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

// Resp RESP protocol reader
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readline reads a line from the RESP stream
func (r *Resp) readline() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n++
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

// readInteger reads the integer value from the RESP stream
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readline()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return int(i64), n, nil
}

// readArray reads the array value from the RESP stream
// *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"
	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	v.array = make([]Value, 0)
	for i := 0; i < length; i++ {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		v.array = append(v.array, v)
	}
	return v, nil
}

// readBulk reads the bulk value from the RESP stream
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"
	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	bulk := make([]byte, length)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return Value{}, err
	}
	v.bulk = string(bulk)
	// 消费掉\r\n
	_, _, err = r.readline()
	if err != nil {
		return Value{}, err
	}
	return v, nil
}
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}
