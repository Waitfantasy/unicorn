package id

const (
	MachineBits         = 10
	SecondSeqBits       = 20
	SecondTimestampBits = 30
	MilliSecondSeqBits  = 10
	MilliTimestampBits  = 40
	ServiceBits         = 1
	IdTypeBits          = 1
	VersionBits         = 1
)

type MetaData struct {
	epoch     uint64
	seq       int
	machineId int
	timestamp int
	version   int
	service   int
	idType    int
}

type Meta struct {
	MachineBit   int8
	SeqBit       int8
	TimestampBit int8
	IdTypeBit    int8
	VersionBit   int8
	data         *MetaData
}

func NewMeta1(machineBit, seqBit, timestampBit, idTypeBit, versionBit int8) *Meta {
	return &Meta{
		MachineBit:   machineBit,
		SeqBit:       seqBit,
		TimestampBit: timestampBit,
		IdTypeBit:    idTypeBit,
		VersionBit:   versionBit,
	}
}

func NewMeta(data *MetaData) *Meta {
	return &Meta{data: data}
}


func (m *Meta) GetSeqBits() uint64 {
	switch m.data.idType {
	case SecondIdType:
		return SecondSeqBits
	case MilliSecondIdType:
		return MilliSecondSeqBits
	default:
		return SecondSeqBits
	}
}

func (m *Meta) GetTimestampBits() uint64 {
	switch m.data.idType {
	case SecondIdType:
		return SecondTimestampBits
	case MilliSecondIdType:
		return MilliTimestampBits
	default:
		return SecondTimestampBits
	}
}

func (m *Meta) GetSeqShift() uint64{
	return uint64(m.MachineBit)
}

func (m *Meta) GetSequenceLeftShift() uint64 {
	return MachineBits
}

func (m *Meta) GetTimestampShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit)
}

func (m *Meta) GetTimestampLeftShift() uint64 {
	return MachineBits + m.GetSeqBits()
}


func (m *Meta) GetServiceLeftShift() uint64 {
	return MachineBits + m.GetSeqBits() + m.GetTimestampBits()
}

func (m *Meta) GetIdTypeShift() uint64 {
	return uint64(m.MachineBit + m.SeqBit + m.TimestampBit)
}

func (m *Meta) GetIdTypeLeftShift() uint64 {
	return MachineBits + m.GetSeqBits() + m.GetTimestampBits() + ServiceBits
}

func (m *Meta) GetVersionShift() uint64{
	return uint64(m.MachineBit + m.SeqBit + m.TimestampBit + m.VersionBit)
}

func (m *Meta) GetVersionLeftShift() uint64 {
	return MachineBits + m.GetSeqBits() + m.GetTimestampBits() + ServiceBits + IdTypeBits
}

func (m *Meta) GetMaxMachine() int64 {
	return -1 ^ (- 1 << MachineBits)
}

func (m *Meta) GetMaxSequence() int64 {
	switch m.data.idType {
	case SecondIdType:
		return -1 ^ (-1 << SecondSeqBits)
	case MilliSecondIdType:
		return -1 ^ (-1 << MilliSecondSeqBits)
	default:
		return -1 ^ (-1 << SecondSeqBits)
	}
}

func (m *Meta) GetMaxTimestamp() int64 {
	return -1 ^ -1<<m.GetTimestampBits()
}

func (m *Meta) GetMaxService() int64 {
	return -1 ^ -1<<ServiceBits
}

func (m *Meta) GetMaxIdType() int64 {
	return -1 ^ -1<<IdTypeBits
}

func (m *Meta) GetMaxVersion() int64 {
	return -1 ^ -1<<VersionBits
}
