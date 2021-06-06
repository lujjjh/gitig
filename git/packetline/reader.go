package packetline

import (
	"fmt"
	"io"
)

type Reader struct {
	r   io.Reader
	err error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (r *Reader) read(data []byte) {
	if r.err == nil {
		_, r.err = io.ReadFull(r.r, data)
	}
}

func (r *Reader) ReadPacket() (Packet, error) {
	var header [packetHeaderSize]byte
	r.read(header[:])
	if r.err != nil {
		return nil, r.err
	}
	// TODO: performance?
	hexValue := func(a byte) int {
		if '0' <= a && a <= '9' {
			return int(a - '0')
		}
		if 'a' <= a && a <= 'f' {
			return int(a - 'a' + 10)
		}
		return -1
	}
	hexToByte := func(data []byte) int {
		_ = data[1]
		value := hexValue(data[0])
		if value < 0 {
			return value
		}
		return (value << 4) | hexValue(data[1])
	}
	// endianness?
	size := hexToByte(header[:])
	if size < 0 {
		return nil, fmt.Errorf("invalid packet header: %s", string(header[:]))
	}
	size = (size << 8) | hexToByte(header[2:])
	if size > MaxPacketLineLength {
		return nil, fmt.Errorf("max packet length exceeded: %d", size)
	}
	switch size {
	case 0:
		return flushPacket{}, nil
	case 1:
		return delimiterPacket{}, nil
	case 2:
		return responseEndPacket{}, nil
	default:
		data := make([]byte, size-packetHeaderSize)
		r.read(data)
		if r.err != nil {
			return nil, r.err
		}
		return dataPacket{data: data}, nil
	}
}
