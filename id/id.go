package id

type ExtractData struct {
	MachineId int
	Sequence  uint64
	Timestamp uint64
	Reserved  int
	IdType    int
	Version   int
}

type Id struct {
	machineId int
	idType    int
	reserved  int
	version   int
	Meta      *Meta
	TimerUtil *TimerUtil
}

func NewId(machineId, idType, version int, epoch uint64) *Id {
	id := &Id{
		machineId: machineId,
		idType:    idType,
		version:   version,
		TimerUtil: NewTimerUtil(idType, epoch),
	}

	switch idType {
	case SecondIdType:
		id.Meta = SecondMeta
	case MilliSecondIdType:
		id.Meta = MilliSecondMeta
	default:
		id.Meta = SecondMeta
	}
	return id
}

func (id *Id) calculate(sequence, lastTimestamp uint64) uint64 {
	var uuid uint64
	uuid |= uint64(id.machineId)
	uuid |= uint64(sequence << id.Meta.GetSeqShift())
	uuid |= uint64(lastTimestamp << id.Meta.GetTimestampShift())
	uuid |= uint64(id.reserved << id.Meta.GetReservedShift())
	uuid |= uint64(id.idType << id.Meta.GetIdTypeShift())
	uuid |= uint64(id.version << id.Meta.GetVersionShift())
	return uuid
}

func (id *Id) transfer(uuid uint64) *ExtractData {
	data := &ExtractData{}
	data.MachineId = int(uuid & uint64(id.Meta.GetMaxMachine()))
	data.Sequence = (uuid >> id.Meta.GetSeqShift()) & uint64(id.Meta.GetMaxSequence())
	data.Timestamp = id.TimerUtil.ConvertTimestamp(uuid >> id.Meta.GetTimestampShift() & uint64(id.Meta.GetMaxTimestamp()))
	data.Reserved = int((uuid >> id.Meta.GetReservedShift()) & uint64(id.Meta.GetMaxReserved()))
	data.IdType = int((uuid >> id.Meta.GetIdTypeShift()) & uint64(id.Meta.GetMaxIdType()))
	data.Version = int((uuid >> id.Meta.GetVersionShift()) & uint64(id.Meta.GetMaxVersion()))
	return data
}
