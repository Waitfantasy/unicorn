package restful

type Config struct {
	Addr       string
	EnableTLS  bool
	CaFile     string
	CertFile   string
	KeyFile    string
	ClientAuth bool
}
