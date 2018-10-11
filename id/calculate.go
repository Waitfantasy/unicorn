package id

func Calculate(sequence uint64, timestamp uint64, meta *Meta) uint64 {
	var uuid uint64
	uuid |= uint64(meta.config.Version << meta.GetVersionLeftShift())
	uuid |= uint64(meta.config.IdGenType << meta.GetIdTypeLeftShift())
	uuid |= uint64(meta.config.ReleaseType << meta.GetReleaseTypeLeftShift())
	uuid |= uint64(timestamp << meta.GetTimestampLeftShift())
	uuid |= uint64(sequence << meta.GetSequenceLeftShift())
	uuid |= uint64(meta.config.MachineId)
	return uuid
}
