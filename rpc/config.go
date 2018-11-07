package rpc

type Config struct {
	Addr       string
	EnableTLS  bool
	CertFile   string
	KeyFile    string
	ServerName string
}

