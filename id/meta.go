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

<<<<<<< HEAD
type MetaData struct {
	epoch     uint64
	seq       int
	machineId int
	timestamp int
	version   int
	service   int
	idType    int
=======
type Meta struct {
	version     int
	epoch       uint64
	idType      int
	releaseType int
	machineId   int
	sequence    uint64
	timestamp   uint64
	config      *IdConfig
>>>>>>> 1be3d383e04f994d82f2912dda4f0a0c3322ce8e
}

type Meta struct {
	data *MetaData
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

func (m *Meta) GetSequenceLeftShift() uint64 {
	return MachineBits
}

func (m *Meta) GetTimestampLeftShift() uint64 {
	return MachineBits + m.GetSeqBits()
}

func (m *Meta) GetServiceLeftShift() uint64 {
	return MachineBits + m.GetSeqBits() + m.GetTimestampBits()
}

func (m *Meta) GetIdTypeLeftShift() uint64 {
	return MachineBits + m.GetSeqBits() + m.GetTimestampBits() + ServiceBits
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
