package master

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/astaxie/beego"
	"github.com/coreos/etcd/clientv3"
	"github.com/sinksmell/bee-crontab/models/common"
)

var (
	// MJobManager master 任务管理器 单例全局变量
	MJobManager *JobMgr
)

// MasterJobMgr  任务管理器
type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

func init() {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
		err    error
	)
	url := beego.AppConfig.String("etcdURL")
	config = clientv3.Config{
		Endpoints:   []string{url},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		beego.Error(err)
		log.Error(err)
		return
	}
	// 得到kv 和lease
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 组装单例
	MJobManager = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

}

// SaveJob 添加或者修改一个任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {

	var (
		jobKey    string
		bytes     []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)
	// 得到job保存路径
	jobKey = getJobKey(job.ID)
	log.Info(jobKey)
	if bytes, err = json.Marshal(job); err != nil {
		return
	}
	// etcd put 操作
	if putResp, err = MJobManager.kv.Put(context.TODO(), jobKey, string(bytes), clientv3.WithPrevKV()); err != nil {
		log.Error(err)
		return
	}
	// 如果prevKV 不为空则返回旧值
	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			// 为了提高容错性
			// 旧值是否正确解析 不影响最终结果
			err = nil
			return
		}
	}
	// 赋值旧值
	oldJob = &oldJobObj
	log.Infof("old job info is: %+v\n", oldJob)

	return
}

// DeleteJob 删除一个任务
func (jobMgr *JobMgr) DeleteJob(job *common.Job) (oldJob *common.Job, err error) {

	var (
		jobKey    string
		oldJobObj common.Job
		delResp   *clientv3.DeleteResponse
	)

	jobKey = getJobKey(job.ID)
	if delResp, err = MJobManager.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}
	// 解析原来的旧值
	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &oldJobObj); err != nil {
			// 是否成功解析出来对操作结果没有影响
			err = nil
			return
		}
	}
	oldJob = &oldJobObj
	return
}

// ListJobs 获取所有的任务
func (jobMgr *JobMgr) ListJobs() (jobs []*common.Job, err error) {
	var (
		allJobKey string
		job       *common.Job
		getResp   *clientv3.GetResponse
	)

	allJobKey = common.JobSavePath
	jobs = make([]*common.Job, 0)
	if getResp, err = jobMgr.kv.Get(context.TODO(), allJobKey, clientv3.WithPrefix()); err != nil {
		return
	}

	if len(getResp.Kvs) != 0 {
		for _, kvPair := range getResp.Kvs {
			job = &common.Job{}
			if err = json.Unmarshal(kvPair.Value, job); err != nil {
				// 容忍了个别任务反序列化失败
				// 正常情况下是可以反序列化的
				err = nil
				continue
			}
			jobs = append(jobs, job)
		}
	}
	return
}

// KillJob  杀死一个任务 向 /cron/killer/JobName put 一个值 worker监听变化,强行终止对应的任务
func (jobMgr *JobMgr) KillJob(job *common.Job) (err error) {

	var (
		leaseID    clientv3.LeaseID
		grantResp  *clientv3.LeaseGrantResponse
		killJobKey = getKillerKey(job.ID)
	)

	// 申请一个租约 设置对应的过期时间
	if grantResp, err = jobMgr.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	leaseID = grantResp.ID
	// 向 /cron/killer/JobName put "kill" 表示杀死对应的任务
	// 租约到期自动删除对应的 k-v
	if _, err = jobMgr.kv.Put(context.TODO(), killJobKey, "kill", clientv3.WithLease(leaseID)); err != nil {
		return
	}
	log.Info("job id", job.ID, "is be killed")

	return
}

// 获取任务存储key
func getJobKey(id string) string {
	return common.JobSavePath + id
}

// 获取要中止任务的key
func getKillerKey(id string) string {
	return common.JobKillerPath + id
}
