package worker

import (
	"context"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// WorkerConfig  Worker节点的配置结构
type Config struct {
	EtcdEndponits   []string `json:"etcd_endponits" yaml:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout" yaml:"etcd_dail_timeout"`
	MongoURL        string   `json:"mongo_url" yaml:"mongo_url"`
	PromPort        int      `json:"prom_port" yaml:"prom_port"`
}

var (
	// WorkerConf Worker的全局配置单例
	Conf *Config
)

// InitConfig 解析Worker配置文件
func InitConfig(ctx context.Context, filename string) (err error) {

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	var (
		content []byte
		config  Config
	)

	// 读取文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		log.Errorf("read config err %w", err)
		return
	}
	// 解析json
	if err = yaml.Unmarshal(content, &config); err != nil {
		log.Errorf("unmarshal config err %w", err)
		return
	}
	Conf = &config
	log.Infof("config %+v\n", Conf)

	return
}
