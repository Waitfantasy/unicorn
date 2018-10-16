package id

const (
	machineBits         = 10
	secondSeqBits       = 20
	secondTimestampBits = 30
	milliSecondSeqBits  = 10
	milliTimestampBits  = 40
	reservedBits        = 2
	idTypeBits          = 1
	versionBits         = 1
)

type meta struct {
	machineBit   uint8
	seqBit       uint8
	timestampBit uint8
	reservedBit  uint8
	idTypeBit    uint8
	versionBit   uint8
}

var secondMeta = &meta{
	machineBit:   machineBits,
	seqBit:       secondSeqBits,
	timestampBit: secondTimestampBits,
	reservedBit:  reservedBits,
	idTypeBit:    idTypeBits,
	versionBit:   versionBits,
}

var milliSecondMeta = &meta{
	machineBit:   machineBits,
	seqBit:       milliSecondSeqBits,
	timestampBit: milliTimestampBits,
	reservedBit:  reservedBits,
	idTypeBit:    idTypeBits,
	versionBit:   versionBits,
}

func (m *meta) GetSeqShift() uint64 {
	return uint64(m.machineBit)
}

func (m *meta) GetTimestampShift() uint64 {
	return uint64(m.machineBit + m.seqBit)
}

func (m *meta) GetReservedShift() uint64 {
	return uint64(m.machineBit + m.seqBit + m.timestampBit)
}

func (m *meta) GetIdTypeShift() uint64 {
	return uint64(m.machineBit + m.seqBit + m.timestampBit + m.reservedBit)
}

func (m *meta) GetVersionShift() uint64 {
	return uint64(m.machineBit + m.seqBit + m.timestampBit + m.reservedBit + m.idTypeBit)
}

func (m *meta) GetMaxMachine() int64 {
	return -1 ^ - 1<<m.machineBit
}

func (m *meta) GetMaxSequence() int64 {
	return -1 ^ -1<<m.seqBit
}

func (m *meta) GetMaxTimestamp() int64 {
	return -1 ^ -1<<m.timestampBit
}

func (m *meta) GetMaxReserved() int64 {
	return -1 ^ -1<<m.reservedBit
}

func (m *meta) GetMaxIdType() int64 {
	return -1 ^ -1<<m.idTypeBit
}

func (m *meta) GetMaxVersion() int64 {
	return -1 ^ -1<<m.versionBit
}
