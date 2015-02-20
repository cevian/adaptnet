package udp

type Report struct {
	AvgLatencyMs int64 //needs to be signed, negative skew clock
	NumBytes     uint32
	Packets      uint32
}
