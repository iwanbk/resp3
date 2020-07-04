package resp3

import (
	"bufio"
	"fmt"
	"io"
	"strconv" // for converting integers to strings
)

// Writer is a redis client writer.
// he RESP3 protocol can be used asymmetrically, as it is in Redis:
// only a subset can be sent by the client to the server,
// while the server can return the full set of types available.
// This is due to the fact that RESP is designed to send non structured commands like SET mykey somevalue or SADD myset a b c d.
// Such commands can be represented as arrays, where each argument is an array element,
// so this is the only type the client needs to send to a server.
type Writer struct {
	*bufio.Writer
}

// NewWriter returns a redis client writer.
func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		Writer: bufio.NewWriter(writer),
	}
}

func (w *Writer) SendCommands(args ...interface{}) error {
	// write the array flag
	w.WriteByte(TypeArray)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(CRLFByte)

	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			w.writeString(v)
		case int64:
			w.writeString(strconv.FormatInt(v, 10))
		case int:
			w.writeString(strconv.FormatInt(int64(v), 10))
		case []byte:
			w.writeString(string(v))
		case bool:
			w.writeString(strconv.FormatBool(v))
		case float32:
			w.writeString(strconv.FormatFloat(float64(v), 'g', -1, 64))
		case float64:
			w.writeString(strconv.FormatFloat(v, 'g', -1, 64))
		default:
			return fmt.Errorf("unsupported data type: %v", arg)
		}
	}
	return w.Flush()
}

func (w *Writer) writeString(str string) {
	w.WriteByte(TypeBlobString)
	w.WriteString(strconv.Itoa(len(str)))
	w.Write(CRLFByte)
	w.WriteString(str)
	w.Write(CRLFByte)
}

// WriteCommand writes a redis command.
func (w *Writer) WriteCommand(args ...string) (err error) {
	// write the array flag
	w.WriteByte(TypeArray)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(CRLFByte)
	// write blobstring
	for _, arg := range args {
		w.WriteByte(TypeBlobString)
		w.WriteString(strconv.Itoa(len(arg)))
		w.Write(CRLFByte)
		w.WriteString(arg)
		w.Write(CRLFByte)
	}
	return w.Flush()
}

// WriteByteCommand writes a redis command in bytes.
func (w *Writer) WriteByteCommand(args ...[]byte) (err error) {
	// write the array flag
	w.WriteByte(TypeArray)
	w.WriteString(strconv.Itoa(len(args)))
	w.Write(CRLFByte)
	// write blobstring
	for _, arg := range args {
		w.WriteByte(TypeBlobString)
		w.WriteString(strconv.Itoa(len(arg)))
		w.Write(CRLFByte)
		w.Write(arg)
		w.Write(CRLFByte)
	}
	return w.Flush()
}
