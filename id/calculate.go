package id

func Calculate(sequence uint64, timestamp uint64, meta *Meta) uint64 {
	var uuid uint64
	uuid |= uint64(meta.data.machineId)
	uuid |= uint64(sequence << meta.GetSequenceLeftShift())
	uuid |= uint64(timestamp << meta.GetTimestampLeftShift())
	uuid |= uint64(meta.data.service << meta.GetServiceLeftShift())
	uuid |= uint64(meta.data.idType << meta.GetIdTypeLeftShift())
	uuid |= uint64(meta.data.version << meta.GetVersionLeftShift())
	return uuid
}
