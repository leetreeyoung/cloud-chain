package worker

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/sinksmell/bee-crontab/models/common"
)

// Register 服务注册 注册节点到etcd
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
	info   string // 本机ip+pid信息
}

var (
	// WorkerRegister worker节点服务注册单例
	WorkerRegister *Register
)

// InitRegister  初始化服务注册单例
func InitRegister() (err error) {

	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
		ip     string
		pid    int
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   Conf.EtcdEndponits,
		DialTimeout: time.Duration(Conf.EtcdDialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		log.Error(err)
		return
	}

	// 创建kv lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	if ip, err = GetLocalIP(); err != nil {
		log.Error(err)
		return
	}

	pid = os.Getpid()

	// 初始化单例
	WorkerRegister = &Register{
		client: client,
		kv:     kv,
		lease:  lease,
		info:   fmt.Sprintf("ip: %s,pid: %d", ip, pid),
	}
	log.Info("start register worker node")
	go WorkerRegister.keepOnline()

	return
}

// 注册ip到 /cron/workers/IP 并自动续租
func (register *Register) keepOnline() {
	var (
		key          string
		leaseResp    *clientv3.LeaseGrantResponse
		leaseID      clientv3.LeaseID
		ctx          context.Context
		cancelFunc   context.CancelFunc
		keepRespChan <-chan *clientv3.LeaseKeepAliveResponse
		err          error
	)

	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		//  续租应答协程
		// 当通道被关闭时 程序协程自动退出
		for keepResp = range keepRespChan {
			if keepResp == nil {
				return
			}
		}

	}()

	//拼接 etcd 中的key 服务注册key
	key = common.JobWorkerPath + register.info

	for {
		// 初始化上下文取消函数
		ctx, cancelFunc = context.WithCancel(context.TODO())
		// 创建租约
		if leaseResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			// 一段时间后重新尝试创建租约
			goto RETRY
		}

		// 自动续租
		leaseID = leaseResp.ID
		if keepRespChan, err = register.lease.KeepAlive(ctx, leaseID); err != nil {
			goto RETRY
		}

		// 注册到etcd
		if _, err = register.kv.Put(ctx, key, "running", clientv3.WithLease(leaseID)); err != nil {
			goto RETRY
		}

	RETRY:
		time.Sleep(time.Second)
		// 取消租约
		cancelFunc()
	}

}

// GetLocalIP 获取本地ip
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs []net.Addr
		addr  net.Addr
		ipNet *net.IPNet
		ok    bool
	)

	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	// 取第一个非lo的网卡
	for _, addr = range addrs {
		// addr是一个接口
		// 使用类型断言
		// 判断是否为ip地址 有可能是unix socket
		if ipNet, ok = addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			// 只接受ipv4
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = common.ErrNoIpFound
	return
}
