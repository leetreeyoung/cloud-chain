package worker

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/sinksmell/bee-crontab/models/common"
)

// JobManager worker 任务管理器
type JobManager struct {
	client  *clientv3.Client //连接etcd 客户端
	kv      clientv3.KV      //kv
	lease   clientv3.Lease   //租约
	watcher clientv3.Watcher //监听器
}

var (
	// Manager Worker全局任务管理器
	Manager *JobManager
)

// InitJobMgr 初始化worker全局任务管理器
func InitJobMgr() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   Conf.EtcdEndponits,
		DialTimeout: time.Duration(Conf.EtcdDialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// 创建kv lease watcher
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	watcher = clientv3.NewWatcher(client)

	// 初始化单例
	Manager = &JobManager{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	// 启动监听任务
	if err = Manager.WatchJobs(); err != nil {
		log.Errorf("watch jobs err:%v", err)
		return
	}

	// 启动监听killer
	if err = Manager.WatchKillers(); err != nil {
		return
	}

	return
}

// WatchJobs 从etcd中读取任务 监听kv变化
func (jobMgr *JobManager) WatchJobs() (err error) {

	var (
		getResp           *clientv3.GetResponse
		kvPair            *mvccpb.KeyValue
		job               *common.Job
		jobID             string
		watchStartRevison int64
		watchChan         clientv3.WatchChan
		watchResp         clientv3.WatchResponse
		watchEvent        *clientv3.Event
		jobEvent          *common.JobEvent
	)

	// 1.get /cron/jobs/下所有任务 并获取 revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JobSavePath, clientv3.WithPrefix()); err != nil {
		log.Errorf("get all jobs err:%v\n", err)
		return
	}

	// 遍历kv
	for _, kvPair = range getResp.Kvs {
		// 反序列化 value->job
		// 如果某次失败则跳过 提高容错性
		if job, err = common.UnpackJob(kvPair.Value); err == nil {
			// 说明是有效的任务
			// 发送给调度器
			//TODO:构造事件 发送给调度器
			jobEvent = common.NewJobEvent(common.JobEventSave, job)
			BeeScheduler.PushJobEvent(jobEvent)
			// log.Infof("send job event :%+v\n", jobEvent)
		}
	}

	//2.从当前revision之后监听变化
	go func() {
		//监听协程
		watchStartRevison = getResp.Header.Revision + 1
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JobSavePath, clientv3.WithPrefix())
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					// 任务保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						// 任务解析失败 跳过
						continue
					}
					// 构造一个更新事件
					jobEvent = common.NewJobEvent(common.JobEventSave, job)
					//  传给调度器
					BeeScheduler.PushJobEvent(jobEvent)
					log.Infof("send put job event :%+v\n", jobEvent)
				case mvccpb.DELETE:
					// 任务删除事件
					jobID = common.ExtractJobID(string(watchEvent.Kv.Key))
					job = &common.Job{ID: jobID}
					// 构造一个任务删除事件
					jobEvent = common.NewJobEvent(common.JobEventDelete, job)
					// 推送给调度器
					BeeScheduler.PushJobEvent(jobEvent)
					log.Infof("send delete job event :%+v\n", jobEvent)
				}
			}
		}

	}()

	return
}

// WatchKillers 从etcd读取killer 监听kv变化
func (jobMgr *JobManager) WatchKillers() (err error) {

	// 监听 /cron/killer/ 目录的变化
	var (
		getResp           *clientv3.GetResponse
		watchChan         clientv3.WatchChan
		watchResp         clientv3.WatchResponse
		watchEvent        *clientv3.Event
		jobEvent          *common.JobEvent
		jobID             string
		job               *common.Job
		watchStartRevison int64
	)

	// 1.get /cron/jobs/ 下的所有任务,并获取当前revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JobKillerPath, clientv3.WithPrefix()); err != nil {
		return
	}
	//从最新的revision之后监听变化
	go func() {
		watchStartRevison = getResp.Header.Revision
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JobKillerPath, clientv3.WithPrefix(), clientv3.WithRev(watchStartRevison))
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					// 杀死某个任务
					// 从key中提取出任务名
					jobID = common.ExtractKillerID(string(watchEvent.Kv.Key))
					job = &common.Job{ID: jobID}
					jobEvent = common.NewJobEvent(common.JobEventKill, job)
					// 事件推送给 schedular
					BeeScheduler.PushJobEvent(jobEvent)
					log.Infof("send kill job event : %v\n", jobEvent)
				case mvccpb.DELETE:
					// killer 任务过期
				}
				// 变化推送给 scheduler
				// scheduler得知后调用cancelFunc取消对应的任务执行
			}
		}
	}()

	return
}

// NewLock 创建分布式锁
func (jobMgr *JobManager) NewLock(jobID string) (lock *JobLock) {
	return InitJobLock(jobID, jobMgr.kv, jobMgr.lease)
}
