package worker

type Config struct {
	EtcdEndpoints   string `json:"etcd_endpoints"`
	EtcdDialTimeout int    `json:"etcd_dial_timeout"`
}

var (
	G_config *Config
)


