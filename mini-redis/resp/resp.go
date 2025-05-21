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
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

// Marshal converts the Value struct to bytes
func (v Value) Marshal() []byte {
	switch v.Typ {
	case "string":
		return v.marshalString()
	case "bulk":
		return v.marshalBulk()
	case "array":
		return v.marshalArray()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

// marshalString converts simple string to bytes
func (v Value) marshalString() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalBulk converts Bulk string to bytes
func (v Value) marshalBulk() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, []byte(v.Bulk)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

// marshalArray converts Array to bytes
func (v Value) marshalArray() []byte {
	length := len(v.Array)
	bytes := make([]byte, 0)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.Array))...)
	bytes = append(bytes, '\r', '\n')
	for i := 0; i < length; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}
	return bytes
}

func (v Value) marshalError() []byte {
	bytes := make([]byte, 0)
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
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

// readArray reads the Array value from the RESP stream
func (r *Resp) readArray() (Value, error) {
	av := Value{}
	av.Typ = "array"
	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	av.Array = make([]Value, 0)
	for i := 0; i < length; i++ {
		v, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		av.Array = append(av.Array, v)
	}
	return av, nil
}

// readBulk reads the Bulk value from the RESP stream
func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.Typ = "bulk"
	length, _, err := r.readInteger()
	if err != nil {
		return Value{}, err
	}
	bulk := make([]byte, length)
	_, err = r.reader.Read(bulk)
	if err != nil {
		return Value{}, err
	}
	v.Bulk = string(bulk)
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

// Writer is a RESP writer
type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// Write is a function that writes takes  Value and writes the types from marshal method
func (w *Writer) Write(v Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		fmt.Println("write error:", err)
		return err
	}
	return nil
}
