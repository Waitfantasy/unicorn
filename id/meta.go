package id

const (
	MachineBits         = 10
	SecondSeqBits       = 20
	SecondTimestampBits = 30
	MilliSecondSeqBits  = 10
	MilliTimestampBits  = 40
	ReservedBits        = 2
	IdTypeBits          = 1
	VersionBits         = 1
)

type Meta struct {
	MachineBit   uint8
	SeqBit       uint8
	TimestampBit uint8
	ReservedBit  uint8
	IdTypeBit    uint8
	VersionBit   uint8
}

var SecondMeta = &Meta{
	MachineBit:   MachineBits,
	SeqBit:       SecondSeqBits,
	TimestampBit: SecondTimestampBits,
	ReservedBit:  ReservedBits,
	IdTypeBit:    IdTypeBits,
	VersionBit:   VersionBits,
}

var MilliSecondMeta = &Meta{
	MachineBit:   MachineBits,
	SeqBit:       MilliSecondSeqBits,
	TimestampBit: MilliTimestampBits,
	ReservedBit:  ReservedBits,
	IdTypeBit:    IdTypeBits,
	VersionBit:   VersionBits,
}

func (m *Meta) GetSeqShift() uint64 {
	return uint64(m.MachineBit)
}

func (m *Meta) GetTimestampShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit)
}

func (m *Meta) GetReservedShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit + m.TimestampBit)
}

func (m *Meta) GetIdTypeShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit + m.TimestampBit + m.ReservedBit)
}

func (m *Meta) GetVersionShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit + m.TimestampBit + m.ReservedBit + m.IdTypeBit)
}

func (m *Meta) GetMaxMachine() int64 {
	return -1 ^ - 1<<m.MachineBit
}

func (m *Meta) GetMaxSequence() int64 {
	return -1 ^ -1<<m.SeqBit
}

func (m *Meta) GetMaxTimestamp() int64 {
	return -1 ^ -1<<m.TimestampBit
}

func (m *Meta) GetMaxReserved() int64 {
	return -1 ^ -1<<m.ReservedBit
}

func (m *Meta) GetMaxIdType() int64 {
	return -1 ^ -1<<m.IdTypeBit
}

func (m *Meta) GetMaxVersion() int64 {
	return -1 ^ -1<<m.VersionBit
}
