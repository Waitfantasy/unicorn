package id

import (
	"fmt"
)

type ExtractData struct {
	MachineId int
	Sequence  uint64
	Timestamp uint64
	Reserved  int
	IdType    int
	Version   int
}

type Config struct {
	Epoch     uint64
	MachineId int
	IdType    int
	Reserved  int
	Version   int
}

type Id struct {
	cfg       *Config
	meta      *meta
	timerUtil *TimerUtil
}

func NewId(cfg *Config) (*Id, error) {
	id := &Id{
		cfg: cfg,
		timerUtil: &TimerUtil{
			cfg.IdType, cfg.Epoch,
		},
	}

	if err := checkConfig(cfg); err != nil {
		return nil, err
	}

	switch id.cfg.IdType {
	case SecondIdType:
		id.meta = secondMeta
	case MilliSecondIdType:
		id.meta = milliSecondMeta
	default:
		id.meta = secondMeta
	}
	return id, nil
}

func checkConfig(cfg *Config) error {
	if cfg.Epoch == 0 {
		return fmt.Errorf("epoch cannot be empty, the id type supports: : \n\t%d: max peak type\n\t%d: min granularity type\n",
			SecondIdType, MilliSecondIdType)
	}

	if cfg.MachineId < 1 || cfg.MachineId > 1024 {
		return fmt.Errorf("machine id is not in range, machine id range: %d ~ %d\n",
			1, 1024)
	}

	if cfg.Version != UnavailableVersion && cfg.Version != NormalVersion {
		return fmt.Errorf("version is unsupported value, the version supports: : \n\t%d: unavailable version\n\t%d: normal version\n",
			UnavailableVersion, NormalVersion)
	}

	return nil
}

func (id *Id) calculate(sequence, lastTimestamp uint64) uint64 {
	var uuid uint64
	uuid |= uint64(id.cfg.MachineId)
	uuid |= uint64(sequence << id.meta.GetSeqShift())
	uuid |= uint64(lastTimestamp << id.meta.GetTimestampShift())
	uuid |= uint64(id.cfg.Reserved << id.meta.GetReservedShift())
	uuid |= uint64(id.cfg.IdType << id.meta.GetIdTypeShift())
	uuid |= uint64(id.cfg.Version << id.meta.GetVersionShift())
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
