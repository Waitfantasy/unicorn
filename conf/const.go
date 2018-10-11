package conf

const (
	IdPeakGenType        = 0
	IdGranularityGenType = 1

	// Machine Id Get Type
	MachineIdByLocal = "local"
	MachineIdByEtcd  = "etcd"

	// Releases Type
	ReleaseLocal = 1
	ReleaseHttp  = 2
	ReleaseGRpc  = 3
)
