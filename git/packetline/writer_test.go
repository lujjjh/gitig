package packetline_test

import (
	"bytes"
	"io"
	"runtime/debug"
	"testing"

	"github.com/lujjjh/gitig/git/packetline"
)

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Log(debug.Stack())
		t.Fatalf("no error expected, got: %v", err)
	}
}

func TestWriter(t *testing.T) {
	testCases := []struct {
		expected []byte
		f        func(w *packetline.Writer)
	}{
		{
			expected: []byte("0000"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WriteFlushPacket())
			},
		},
		{
			expected: []byte("0001"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WriteDelimiterPacket())
			},
		},
		{
			expected: []byte("0002"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WriteResponseEndPacket())
			},
		},
		{
			expected: []byte("0006a\n"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WritePacket([]byte("a\n")))
			},
		},
		{
			expected: []byte("0005a"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WritePacket([]byte("a")))
			},
		},
		{
			expected: []byte("000bfoobar\n"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WritePacket([]byte("foobar\n")))
			},
		},
		{
			expected: []byte("0004"),
			f: func(w *packetline.Writer) {
				assertNoError(t, w.WritePacket([]byte("")))
			},
		},
	}
	for i, tc := range testCases {
		var buf bytes.Buffer
		w := packetline.NewWriter(&buf)
		tc.f(w)
		actual := buf.Bytes()
		if !bytes.Equal(tc.expected, actual) {
			t.Errorf("testCase[%d] failed: expected: %s, got: %s", i, string(tc.expected), string(actual))
		}
	}
}

func TestWriterMaxLengthExceeded(t *testing.T) {
	var data [65517]byte
	w := packetline.NewWriter(io.Discard)
	if err := w.WritePacket(data[:]); err != packetline.ErrMaxPacketLineLengthExceeded {
		t.Errorf("expected: ErrMaxPacketLineLengthExceeded, got: %v", err)
	}
}
