package packetline

type Packet interface {
	FlushPacket() bool
	DelimiterPacket() bool
	ResponseEndPacket() bool
	Data() ([]byte, bool)
}

type flushPacket struct{}

func (flushPacket) FlushPacket() bool       { return true }
func (flushPacket) DelimiterPacket() bool   { return false }
func (flushPacket) ResponseEndPacket() bool { return false }
func (flushPacket) Data() ([]byte, bool)    { return nil, false }

type delimiterPacket struct{}

func (delimiterPacket) FlushPacket() bool       { return false }
func (delimiterPacket) DelimiterPacket() bool   { return true }
func (delimiterPacket) ResponseEndPacket() bool { return false }
func (delimiterPacket) Data() ([]byte, bool)    { return nil, false }

type responseEndPacket struct{}

func (responseEndPacket) FlushPacket() bool       { return false }
func (responseEndPacket) DelimiterPacket() bool   { return false }
func (responseEndPacket) ResponseEndPacket() bool { return true }
func (responseEndPacket) Data() ([]byte, bool)    { return nil, false }

type dataPacket struct {
	data []byte
}

func (dataPacket) FlushPacket() bool       { return false }
func (dataPacket) DelimiterPacket() bool   { return false }
func (dataPacket) ResponseEndPacket() bool { return false }
func (d dataPacket) Data() ([]byte, bool)  { return d.data, true }
