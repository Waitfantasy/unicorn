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
	meta      *meta
	timerUtil *TimerUtil
}

func NewId(machineId, idType, version int, epoch uint64) *Id {
	id := &Id{
		machineId: machineId,
		idType:    idType,
		version:   version,
		timerUtil: &TimerUtil{
			idType, epoch,
		},
	}

	switch idType {
	case SecondIdType:
		id.meta = secondMeta
	case MilliSecondIdType:
		id.meta = milliSecondMeta
	default:
		id.meta = secondMeta
	}
	return id
}

func (id *Id) calculate(sequence, lastTimestamp uint64) uint64 {
	var uuid uint64
	uuid |= uint64(id.machineId)
	uuid |= uint64(sequence << id.meta.GetSeqShift())
	uuid |= uint64(lastTimestamp << id.meta.GetTimestampShift())
	uuid |= uint64(id.reserved << id.meta.GetReservedShift())
	uuid |= uint64(id.idType << id.meta.GetIdTypeShift())
	uuid |= uint64(id.version << id.meta.GetVersionShift())
	return uuid
}

func (id *Id) transfer(uuid uint64) *ExtractData {
	data := &ExtractData{}
	data.MachineId = int(uuid & uint64(id.meta.GetMaxMachine()))
	data.Sequence = (uuid >> id.meta.GetSeqShift()) & uint64(id.meta.GetMaxSequence())
	data.Timestamp = id.timerUtil.ConvertTimestamp(uuid >> id.meta.GetTimestampShift() & uint64(id.meta.GetMaxTimestamp()))
	data.Reserved = int((uuid >> id.meta.GetReservedShift()) & uint64(id.meta.GetMaxReserved()))
	data.IdType = int((uuid >> id.meta.GetIdTypeShift()) & uint64(id.meta.GetMaxIdType()))
	data.Version = int((uuid >> id.meta.GetVersionShift()) & uint64(id.meta.GetMaxVersion()))
	return data
}
