package master

import (
	"context"
	log "github.com/sirupsen/logrus"
	"time"

	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sinksmell/bee-crontab/models/common"
)

// WorkerMgr  worker管理 用来发现worker
// /cron/worker/
type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

// WorkerInfo 为了使获取节点列表 得到的信息更加丰富 而不是单纯的ip
// 从而添加的描述节点状态的结构体 方便之后拓展
type WorkerInfo struct {
	Time string `json:"time"` // 查询时间
	IP   string `json:"ip"`   // 节点ip
}

var (
	// WorkerManager master用来查看worker 信息的全局单例
	WorkerManager *WorkerMgr
)

func init() {
	InitWorkerMgr()
}

// InitWorkerMgr  初始化全局单例
func InitWorkerMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)
	url := beego.AppConfig.String("etcdURL")
	config = clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		beego.Error(err)
		return
	}
	// 得到kv 和lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	WorkerManager = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

// ListWorkers 获取worker节点的列表
func (workerMgr *WorkerMgr) ListWorkers() (workers []*WorkerInfo, err error) {

	var (
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		ip      string
		info    *WorkerInfo
	)

	workers = make([]*WorkerInfo, 0)
	// 获取目录下所有的节点 ip
	if getResp, err = workerMgr.kv.Get(context.TODO(), common.JobWorkerPath, clientv3.WithPrefix()); err != nil {
		log.Error(err)
		return
	}

	// 保存结果
	for _, kvPair = range getResp.Kvs {
		ip = common.ExtarctWorkerIP(string(kvPair.Key))
		if len(ip) != 0 {
			info = &WorkerInfo{IP: ip}
			info.Time = time.Now().Format("2006-01-02 15:04:05")
			workers = append(workers, info)
		}
	}

	return
}
