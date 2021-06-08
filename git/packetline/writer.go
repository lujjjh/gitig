package packetline

import (
	"errors"
	"fmt"
	"io"
)

const (
	// MaxPacketLineLength indicates the maximum length of a pkt-line,
	// including the 4 bytes packet header.
	MaxPacketLineLength = 65520

	packetHeaderSize = 4
)

var (
	// ErrMaxPacketLineLengthExceeded indicates that the length of pkt-line exceeds
	// MaxPacketLineLength and must not be sent according to the protocol.
	ErrMaxPacketLineLengthExceeded = errors.New("max packet line length exceeded")

	hexChar = []byte("0123456789abcdef")

	flushPacketHeader       = []byte("0000") // indicates the end of a message
	delimiterPacketHeader   = []byte("0001") // separates sections of a message
	responseEndPacketHeader = []byte("0002") // indicates the end of a response for stateless connections
)

type Writer struct {
	w   io.Writer
	err error
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (w *Writer) write(data []byte) {
	if w.err == nil {
		_, w.err = w.w.Write(data)
	}
}

func (w *Writer) Err() error {
	return w.err
}

func (w *Writer) WriteFlushPacket() *Writer {
	w.write(flushPacketHeader)
	return w
}

func (w *Writer) WriteDelimiterPacket() *Writer {
	w.write(delimiterPacketHeader)
	return w
}

func (w *Writer) WriteResponseEndPacket() *Writer {
	w.write(responseEndPacketHeader)
	return w
}

func (w *Writer) packetHeader(size int) [packetHeaderSize]byte {
	var buf [packetHeaderSize]byte
	hex := func(a int) byte { return hexChar[a&15] }
	// endianness?
	buf[0] = hex(size >> 12)
	buf[1] = hex(size >> 8)
	buf[2] = hex(size >> 4)
	buf[3] = hex(size)
	return buf
}

func (w *Writer) writePacketHeader(size int) {
	header := w.packetHeader(size)
	w.write(header[:])
}

func (w *Writer) WritePacket(data []byte) *Writer {
	size := packetHeaderSize + len(data)
	if size > MaxPacketLineLength {
		w.err = ErrMaxPacketLineLengthExceeded
		return w
	}
	w.writePacketHeader(size)
	w.write(data)
	return w
}

func (w *Writer) WritePacketFmt(format string, a ...interface{}) *Writer {
	// TODO: faster conversion from string to readonly []byte?
	return w.WritePacket([]byte(fmt.Sprintf(format, a...)))
}
