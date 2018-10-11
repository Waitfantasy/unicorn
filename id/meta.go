package id

import "github.com/Waitfantasy/unicorn/conf"

const (
	MachineBits        = 10
	peakSeqBits        = 20
	secondTsBits       = 30
	granularitySeqBits = 10
	milliTsBits        = 40
	ReleaseTypeBits    = 2
	IdTypeBits         = 1
	VersionBits        = 1
)

type Meta struct {
	config *IdConfig
}

func NewMeta(c *IdConfig) *Meta {
	return &Meta{config: c}
}

func (m Meta) GetSeqBits() uint64 {
	switch m.config.IdGenType {
	case conf.IdPeakGenType:
		return peakSeqBits
	case conf.IdGranularityGenType:
		return granularitySeqBits
	default:
		return peakSeqBits
	}
}

func (m Meta) GetTimestampBits() uint64 {
	switch m.config.IdGenType {
	case conf.IdPeakGenType:
		return secondTsBits
	case conf.IdGranularityGenType:
		return milliTsBits
	default:
		return secondTsBits
	}
}

func (m Meta) GetMaxSequence() uint64 {
	switch m.config.IdGenType {
	case conf.IdPeakGenType:
		return -1 ^ (-1 << peakSeqBits)
	case conf.IdGranularityGenType:
		return -1 ^ (-1 << granularitySeqBits)
	default:
		return -1 ^ (-1 << peakSeqBits)
	}
}

func (m Meta) GetSequenceLeftShift() uint64 {
	return MachineBits
}

func (m Meta) GetTimestampLeftShift() uint64 {
	return m.GetTimestampBits() + MachineBits
}

func (m Meta) GetReleaseTypeLeftShift() uint64 {
	return m.GetTimestampBits() + m.GetSeqBits() + MachineBits
}

func (m Meta) GetIdTypeLeftShift() uint64 {
	return ReleaseTypeBits + m.GetTimestampBits() + m.GetSeqBits() + MachineBits
}

func (m *Meta) GetVersionLeftShift() uint64 {
	return IdTypeBits + ReleaseTypeBits + m.GetTimestampBits() + m.GetSeqBits() + MachineBits
}
