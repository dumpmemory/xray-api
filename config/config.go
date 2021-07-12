package config

type XRAY struct {
	Port            int
	Grpc            int
	Protocol        string
	Type            string
	Path            string
	Tls             bool
	CertificateFile string
	KeyFile         string
	Method          string
	Network         string
	Dns             []string
	Block           struct {
		BT      bool
		Ips     []string
		Domains []string
	}
	AutoRestart bool
}
type CONF struct {
	Listen      string
	Key         string
	Xray        XRAY
	Syncfile    string
	Autorestart string
	Debug       bool
}
