package packetline_test

import (
	"bytes"
	"runtime/debug"
	"testing"

	"github.com/lujjjh/gitig/git/packetline"
)

func assert(t *testing.T, value bool) {
	if !value {
		t.Log(string(debug.Stack()))
		t.Error("assertion failure")
	}
}

func TestReader(t *testing.T) {
	testCases := []struct {
		payload []byte
		test    func(p packetline.Packet)
	}{
		{
			[]byte("0000"),
			func(p packetline.Packet) {
				assert(t, p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				_, ok := p.Data()
				assert(t, !ok)
			},
		},
		{
			[]byte("0001"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				_, ok := p.Data()
				assert(t, !ok)
			},
		},
		{
			[]byte("0002"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, p.ResponseEndPacket())
				_, ok := p.Data()
				assert(t, !ok)
			},
		},
		{
			[]byte("0006a\n"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				data, ok := p.Data()
				assert(t, ok)
				assert(t, bytes.Equal(data, []byte("a\n")))
			},
		},
		{
			[]byte("0005a"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				data, ok := p.Data()
				assert(t, ok)
				assert(t, bytes.Equal(data, []byte("a")))
			},
		},
		{
			[]byte("000bfoobar\n"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				data, ok := p.Data()
				assert(t, ok)
				assert(t, bytes.Equal(data, []byte("foobar\n")))
			},
		},
		{
			[]byte("0004\n"),
			func(p packetline.Packet) {
				assert(t, !p.FlushPacket())
				assert(t, !p.DelimiterPacket())
				assert(t, !p.ResponseEndPacket())
				data, ok := p.Data()
				assert(t, ok)
				assert(t, bytes.Equal(data, []byte("")))
			},
		},
	}
	for i, tc := range testCases {
		r := packetline.NewReader(bytes.NewReader(tc.payload))
		p, err := r.ReadPacket()
		if err != nil {
			t.Errorf("testCase[%d]: %s", i, err)
			continue
		}
		tc.test(p)
	}
}
