package master

import (
	"encoding/json"
	"io/ioutil"
)
//配置结构体
type Config struct {
	Port            int    `json:"port"`
	ReadTimeout     int    `json:"read_timeout"`
	WriteTimeout    int    `json:"write_timeout"`
	EtcdEndpoints   string `json:"etcd_endpoints"`
	EtcdDialTimeout int    `json:"etcd_dial_timeout"`
}

var(
	G_config *Config
)

//初始化config
func InitConfig(filename string) error {
	var(
		err error
		content []byte
		conf Config
	)

	content,err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//json反序列化
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return err
	}
	//赋值单例
	G_config = &conf

	return nil
}